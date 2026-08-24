package actions

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Scan collects all regular files in the current directory (non-recursive) into
// a map keyed by filename. Directories, symlinks and non-regular files are
// skipped. Stat calls are executed in a bounded worker pool so that the syscalls
// run concurrently across CPU cores.
func Scan(dir string) (map[string]*File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	type job struct {
		entry fs.DirEntry
	}

	type result struct {
		entry fs.DirEntry
		info  fs.FileInfo
		err   error
	}

	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	if workers > len(entries) {
		workers = len(entries)
	}
	if workers < 1 {
		return map[string]*File{}, nil
	}

	jobs := make(chan job)
	results := make(chan result, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				info, err := j.entry.Info()
				results <- result{entry: j.entry, info: info, err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, e := range entries {
			jobs <- job{entry: e}
		}
		close(jobs)
	}()

	files := make(map[string]*File)
	for r := range results {
		if r.err != nil {
			continue
		}
		if !r.info.Mode().IsRegular() {
			continue
		}
		name := r.entry.Name()
		ext := Extension(name)
		files[name] = &File{
			Name:    name,
			Path:    filepath.Join(dir, name),
			Size:    r.info.Size(),
			ModTime: r.info.ModTime(),
			Ext:     ext,
			Media:   ClassifyMedia(ext),
		}
	}
	return files, nil
}
