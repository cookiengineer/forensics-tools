package actions

import (
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// MediaKind classifies a file for the perception step.
type MediaKind int

const (
	MediaNone MediaKind = iota
	MediaImage
	MediaVideo
)

// File is the unit of comparison for the fs-cleanup pipeline.
// The top-level map is keyed by filename (basename).
type File struct {
	Name       string // basename, also the map key
	Path       string // == Name (non-recursive scan)
	Size       int64
	ModTime    time.Time
	Ext        string // lowercased, without the leading dot
	Media      MediaKind
	Hash       string // content sha256, only computed for same-size candidates
	PHash      string // perception hash, only computed for media candidates
	FrameStart string // video perception hash of the frame one second in
	FrameEnd   string // video perception hash of the frame one second before the end
	Width      int
	Height     int
	Duration   float64

	DuplicateOf string  // name of the preserved file this duplicates
	Similarity  float64 // percent similarity to DuplicateOf (for reporting)
	Perceptual  bool    // matched via perceptual analysis (not an exact byte duplicate)
	Delete      bool    // marked for deletion
}

var imageExtensions = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "webp": true, "gif": true, "bmp": true,
}

var videoExtensions = map[string]bool{
	"mkv": true, "mp4": true, "webm": true, "mov": true, "avi": true,
}

// ClassifyMedia returns the MediaKind for an extension (already lowercased,
// without the leading dot).
func ClassifyMedia(ext string) MediaKind {
	if imageExtensions[ext] {
		return MediaImage
	}
	if videoExtensions[ext] {
		return MediaVideo
	}
	return MediaNone
}

// Extension extracts and lowercases the extension of a filename without the dot.
func Extension(name string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
}

// SimilarityPercent returns how similar two file sizes are, as a percentage in
// the range (0, 100]. It is defined as the ratio of the smaller size to the
// larger size, which is what "99% byte size of the other file" describes.
func SimilarityPercent(a, b int64) float64 {
	if a == 0 && b == 0 {
		return 100
	}
	min, max := a, b
	if b < a {
		min, max = b, a
	}
	if max == 0 {
		return 0
	}
	return (float64(min) / float64(max)) * 100
}

// SizeWithinPercent reports whether two file sizes are within the given
// percentage of each other (e.g. 99 means the smaller is >= 99% of the larger).
func SizeWithinPercent(a, b int64, percent float64) bool {
	if a == b {
		return true
	}
	return SimilarityPercent(a, b) >= percent
}

// SimilarSizeThreshold is the percentage used to consider two files of the same
// extension as near-identical candidates for the perception step.
const SimilarSizeThreshold = 99.0

// copyMarkers are substrings that indicate a filename is a copy, backup or
// derivative of another file. Presence of any marker increases the
// originalityScore (making the file less likely to be the original).
var copyMarkers = []string{
	".bak", ".backup", ".old", ".tmp", ".temp", ".temp~", "~",
	" copy", "copy of", " (copy)", "[copy]", " - copy",
	"-copy", "_copy", "-copie", " copie",
	" (1)", " (2)", " (3)", "(1)", "(2)", "(3)",
	"-1", "-2", "-3", "_1", "_2", "_3",
	"-new", "-old", "-final", "-draft", "-v2", "-v3", "_v2", "_v3",
}

// originalityScore rates how likely a filename is to be a derivative rather
// than the original. A higher score means "more likely a copy". Markers are
// matched case-insensitively.
func originalityScore(name string) int {
	lower := strings.ToLower(name)
	score := 0
	for _, marker := range copyMarkers {
		if strings.Contains(lower, marker) {
			score++
		}
	}
	return score
}

// PreserveOriginal chooses the file to preserve from a group of exact
// duplicates. It prefers the filename that most looks like the original; if
// that cannot be inferred, it falls back to the oldest modtime. Ties are broken
// by alphabetical order for reproducibility.
func PreserveOriginal(files []*File) *File {
	if len(files) == 0 {
		return nil
	}
	if len(files) == 1 {
		return files[0]
	}

	best := files[0]
	bestScore := originalityScore(best.Name)
	unique := true

	for _, f := range files[1:] {
		s := originalityScore(f.Name)
		switch {
		case s < bestScore:
			best, bestScore, unique = f, s, true
		case s == bestScore:
			unique = false
		}
	}

	if unique {
		return best
	}

	// No single filename is clearly the original: fall back to oldest modtime,
	// with an alphabetical tie-break.
	for _, f := range files[1:] {
		switch {
		case f.ModTime.Before(best.ModTime):
			best = f
		case f.ModTime.Equal(best.ModTime) && f.Name < best.Name:
			best = f
		}
	}
	return best
}

// SameSecond reports whether two durations have the same whole-second value,
// giving video encoders slack for sub-second differences.
func SameSecond(a, b float64) bool {
	return int64(math.Round(a)) == int64(math.Round(b))
}

// runParallel fans out fn over items using a bounded pool of workers goroutines.
// fn must be safe to run concurrently; each item is processed exactly once.
func runParallel[T any](items []T, workers int, fn func(T)) {
	if len(items) == 0 {
		return
	}
	if workers < 1 {
		workers = 1
	}
	if workers > len(items) {
		workers = len(items)
	}

	jobs := make(chan T)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for it := range jobs {
				fn(it)
			}
		}()
	}
	for _, it := range items {
		jobs <- it
	}
	close(jobs)
	wg.Wait()
}
