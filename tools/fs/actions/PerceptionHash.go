package actions

import (
	"runtime"
	"sort"

	"fs/helpers/ffmpeg"
	"fs/helpers/pdq"
)

// pair is a candidate near-duplicate relationship between two video files of
// the same extension with ~99% similar sizes.
type pair struct {
	a, b *File
}

// PerceptionHash resolves near-duplicates among media files that were not
// already marked as identical (and hence deleted) by the hashing step.
//
// Images: every file of a given extension is compared to every other via its
// PDQ perceptual hash. Files sharing a hash are duplicates and the
// highest-resolution variant is preserved. No size gate is applied because a
// resize changes byte size far more than 1%.
//
// Videos: only files of the same extension with ~99% similar sizes are
// considered. Durations are compared first; if they differ (in whole seconds),
// the shorter video is a candidate. Otherwise the frame one second in and the
// frame one second before the end are compared. Videos are never declared
// duplicates unless their opening frames match perceptually.
func PerceptionHash(files map[string]*File) {
	byExt := make(map[string][]*File)
	for _, f := range files {
		if f.Delete || f.Media == MediaNone {
			continue
		}
		byExt[f.Ext] = append(byExt[f.Ext], f)
	}

	videosAvailable := ffmpeg.Available()
	if videosAvailable {
		probeVideoDurations(byExt)
	}

	for _, group := range byExt {
		if len(group) < 2 {
			continue
		}
		switch group[0].Media {
		case MediaImage:
			processImages(group)
		case MediaVideo:
			if !videosAvailable {
				continue
			}
			sort.Slice(group, func(i, j int) bool {
				if group[i].Size != group[j].Size {
					return group[i].Size < group[j].Size
				}
				return group[i].Name < group[j].Name
			})
			if pairs := buildVideoPairs(group); len(pairs) > 0 {
				processVideoPairs(pairs)
			}
		}
	}
}

// probeVideoDurations fills in durations for videos, in parallel.
func probeVideoDurations(byExt map[string][]*File) {
	var videos []*File
	for _, g := range byExt {
		if len(g) > 0 && g[0].Media == MediaVideo {
			videos = append(videos, g...)
		}
	}
	if len(videos) == 0 {
		return
	}
	runParallel(videos, runtime.NumCPU(), func(f *File) {
		d, err := ffmpeg.ProbeDuration(f.Path)
		if err == nil {
			f.Duration = d
		}
	})
}

// processImages computes the PDQ hash (and pixel dimensions) of every image in
// parallel, then groups by hash and preserves the highest-resolution variant of
// each group.
func processImages(group []*File) {
	runParallel(group, runtime.NumCPU(), func(f *File) {
		h, w, hgt, err := pdq.HashFile(f.Path)
		if err == nil {
			f.PHash = h
			f.Width = w
			f.Height = hgt
		}
	})

	byHash := make(map[string][]*File)
	for _, f := range group {
		if f.PHash != "" {
			byHash[f.PHash] = append(byHash[f.PHash], f)
		}
	}

	for _, dups := range byHash {
		if len(dups) < 2 {
			continue
		}
		keep := highestResolution(dups)
		for _, f := range dups {
			if f == keep {
				continue
			}
			f.DuplicateOf = keep.Name
			f.Similarity = pdq.Similarity(f.PHash, keep.PHash)
			f.Perceptual = true
			f.Delete = true
		}
	}
}

// highestResolution returns the file with the largest pixel area, breaking ties
// by oldest modtime and then alphabetically for reproducibility.
func highestResolution(files []*File) *File {
	best := files[0]
	for _, f := range files[1:] {
		area, bestArea := f.Width*f.Height, best.Width*best.Height
		if area > bestArea {
			best = f
			continue
		}
		if area == bestArea && (f.ModTime.Before(best.ModTime) ||
			(f.ModTime.Equal(best.ModTime) && f.Name < best.Name)) {
			best = f
		}
	}
	return best
}

// buildVideoPairs produces the near-duplicate candidate pairs from a
// size-sorted group. Because the group is sorted by size ascending, once a
// file's size falls below the similarity threshold of the current anchor we can
// stop.
func buildVideoPairs(group []*File) []pair {
	var pairs []pair
	for i := 0; i < len(group); i++ {
		a := group[i]
		if a.Delete {
			continue
		}
		for j := i + 1; j < len(group); j++ {
			b := group[j]
			if b.Delete {
				continue
			}
			if !SizeWithinPercent(a.Size, b.Size, SimilarSizeThreshold) {
				break
			}
			pairs = append(pairs, pair{a: a, b: b})
		}
	}
	return pairs
}

// processVideoPairs resolves near-duplicates among video candidates. Duration
// alone is never enough to declare two videos identical: the opening frame must
// always match perceptually. Equal-duration videos additionally require the
// closing frame to match.
func processVideoPairs(pairs []pair) {
	var sameSecond, diffSecond []pair
	needStart := make(map[*File]bool)
	needEnd := make(map[*File]bool)

	for _, p := range pairs {
		needStart[p.a] = true
		needStart[p.b] = true
		if SameSecond(p.a.Duration, p.b.Duration) {
			sameSecond = append(sameSecond, p)
			needEnd[p.a] = true
			needEnd[p.b] = true
		} else {
			diffSecond = append(diffSecond, p)
		}
	}

	// Extract and hash the opening frame of every candidate video.
	if len(needStart) > 0 {
		videos := filesFromSet(needStart)
		runParallel(videos, runtime.NumCPU(), func(v *File) {
			h, ok := frameHash(v, startSeek(v.Duration))
			if ok {
				v.FrameStart = h
			}
		})
	}
	// Extract and hash the closing frame only for equal-duration candidates.
	if len(needEnd) > 0 {
		videos := filesFromSet(needEnd)
		runParallel(videos, runtime.NumCPU(), func(v *File) {
			seek, ok := endSeek(v.Duration)
			if !ok {
				return
			}
			h, ok := frameHash(v, seek)
			if ok {
				v.FrameEnd = h
			}
		})
	}

	// Differing durations are only near-identical when the videos share the same
	// opening frame; otherwise they are unrelated and are left untouched.
	for _, p := range diffSecond {
		a, b := p.a, p.b
		if a.Delete || b.Delete {
			continue
		}
		if a.FrameStart == "" || b.FrameStart == "" || a.FrameStart != b.FrameStart {
			continue
		}
		markShorter(p)
	}

	// Equal durations are identical only when both opening and closing frames
	// match perceptually.
	for _, p := range sameSecond {
		a, b := p.a, p.b
		if a.Delete || b.Delete {
			continue
		}
		if a.FrameStart == "" || b.FrameStart == "" || a.FrameStart != b.FrameStart {
			continue
		}
		if a.FrameEnd == "" || b.FrameEnd == "" || a.FrameEnd != b.FrameEnd {
			continue
		}
		markNewerForDeletion(a, b)
	}
}

// frameHash extracts a single frame at the given offset and returns its PDQ
// perceptual hash.
func frameHash(v *File, seek float64) (string, bool) {
	frame, err := ffmpeg.ExtractFrame(v.Path, seek)
	if err != nil {
		return "", false
	}
	h, err := pdq.HashBytes(frame)
	if err != nil {
		return "", false
	}
	return h, true
}

// filesFromSet converts a set of files into a slice.
func filesFromSet(set map[*File]bool) []*File {
	files := make([]*File, 0, len(set))
	for f := range set {
		files = append(files, f)
	}
	return files
}

// markShorter marks the shorter-duration video of a pair for deletion.
func markShorter(p pair) {
	a, b := p.a, p.b
	if a.Delete || b.Delete {
		return
	}
	var shorter, longer *File
	switch {
	case a.Duration < b.Duration:
		shorter, longer = a, b
	case b.Duration < a.Duration:
		shorter, longer = b, a
	default:
		return
	}
	shorter.DuplicateOf = longer.Name
	shorter.Similarity = SimilarityPercent(shorter.Size, longer.Size)
	shorter.Perceptual = true
	shorter.Delete = true
}

// markNewerForDeletion preserves the oldest file of an identical pair and marks
// the newer one for deletion.
func markNewerForDeletion(a, b *File) {
	older, newer := a, b
	if b.ModTime.Before(a.ModTime) || (b.ModTime.Equal(a.ModTime) && b.Name < a.Name) {
		older, newer = b, a
	}
	newer.DuplicateOf = older.Name
	newer.Similarity = SimilarityPercent(a.Size, b.Size)
	newer.Perceptual = true
	newer.Delete = true
}

// startSeek returns the offset of the opening comparison frame. For videos
// shorter than two seconds there is no "one second in" frame, so a midpoint
// frame is used instead.
func startSeek(duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	if duration < 2.0 {
		return duration / 2
	}
	return 1.0
}

// endSeek returns the offset of the closing comparison frame (one second before
// the end). For videos shorter than two seconds there is no distinct closing
// frame, which is reported via ok=false.
func endSeek(duration float64) (seek float64, ok bool) {
	if duration < 2.0 {
		return 0, false
	}
	return duration - 1.0, true
}
