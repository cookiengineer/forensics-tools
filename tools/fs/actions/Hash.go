package actions

import (
	"runtime"
	"sort"
	"sync"

	"fs/helpers/sha256sum"
)

// Hash groups files by byte size, hashes the same-size candidates in a bounded
// worker pool, and resolves exact duplicates. Files sharing a size and a hash
// are treated as identical; within each duplicate group PreserveOriginal selects
// the file to keep and the rest are marked for deletion. Same-size files with
// differing hashes are left untouched for the perception step.
func Hash(files map[string]*File) {
	sizeGroups := make(map[int64][]*File)
	for _, f := range files {
		if f.Delete {
			continue
		}
		sizeGroups[f.Size] = append(sizeGroups[f.Size], f)
	}

	// Collect the candidates that actually need hashing (size groups of 2+).
	var candidates []*File
	for _, group := range sizeGroups {
		if len(group) >= 2 {
			candidates = append(candidates, group...)
		}
	}

	hashFiles(candidates)

	// Within each size group, sub-group by hash and resolve exact duplicates.
	for _, group := range sizeGroups {
		if len(group) < 2 {
			continue
		}
		byHash := make(map[string][]*File)
		for _, f := range group {
			if f.Hash != "" {
				byHash[f.Hash] = append(byHash[f.Hash], f)
			}
		}
		for _, dups := range byHash {
			if len(dups) < 2 {
				continue
			}
			keep := PreserveOriginal(dups)
			for _, f := range dups {
				if f == keep {
					continue
				}
				f.DuplicateOf = keep.Name
				f.Similarity = 100
				f.Delete = true
			}
		}
	}
}

// hashFiles computes the content hash for each candidate file in parallel.
func hashFiles(candidates []*File) {
	if len(candidates) == 0 {
		return
	}

	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	if workers > len(candidates) {
		workers = len(candidates)
	}

	type result struct {
		file *File
		hash string
	}

	jobs := make(chan *File)
	results := make(chan result, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range jobs {
				if hash, err := sha256sum.HashFile(f.Path); err == nil {
					results <- result{file: f, hash: hash}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, f := range candidates {
			jobs <- f
		}
		close(jobs)
	}()

	for r := range results {
		r.file.Hash = r.hash
	}
}

// SortedFiles returns the files in deterministic (alphabetical) order by name.
func SortedFiles(files map[string]*File) []*File {
	list := make([]*File, 0, len(files))
	for _, f := range files {
		list = append(list, f)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}
