package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"fs/actions"
)

func main() {
	confirm := flag.Bool("confirm", false, "delete the listed files")
	flag.Parse()

	files, err := actions.Scan(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	actions.Hash(files)
	actions.PerceptionHash(files)

	deleted := make([]*actions.File, 0)
	for _, f := range files {
		if f.Delete {
			deleted = append(deleted, f)
		}
	}
	// Sort by the preserved original first, then the duplicate, so all
	// duplicates of the same original are grouped together. This ordering is
	// also used for the --confirm deletion run to keep listing and deletion in
	// the same reproducible order.
	sort.Slice(deleted, func(i, j int) bool {
		if deleted[i].DuplicateOf != deleted[j].DuplicateOf {
			return deleted[i].DuplicateOf < deleted[j].DuplicateOf
		}
		return deleted[i].Name < deleted[j].Name
	})

	fmt.Println("Files marked for deletion:")
	if len(deleted) == 0 {
		fmt.Println("  (none)")
	} else {
		printDeletionTable(deleted)
	}

	fmt.Println()
	fmt.Println("Almost identical:")
	similar := 0
	for _, f := range deleted {
		if f.DuplicateOf == "" || !f.Perceptual {
			continue
		}
		fmt.Printf("  %s is %.2f%% identical to %s\n", f.Name, f.Similarity, f.DuplicateOf)
		similar++
	}
	if similar == 0 {
		fmt.Println("  (none)")
	}

	if *confirm && len(deleted) > 0 {
		fmt.Println()
		fmt.Println("Deleting...")
		for _, f := range deleted {
			if err := os.Remove(f.Path); err != nil {
				fmt.Fprintf(os.Stderr, "  error deleting %s: %v\n", f.Name, err)
				continue
			}
			fmt.Printf("  deleted %s\n", f.Name)
		}
	}
}

// printDeletionTable prints the deletion list as a two-column table: the
// preserved original on the left and the duplicate marked for deletion on the
// right, with columns aligned to the widest entry.
func printDeletionTable(deleted []*actions.File) {
	const (
		originalHeader  = "Original"
		duplicateHeader = "Duplicate"
	)
	originalWidth := len(originalHeader)
	duplicateWidth := len(duplicateHeader)
	for _, f := range deleted {
		if len(f.DuplicateOf) > originalWidth {
			originalWidth = len(f.DuplicateOf)
		}
		if len(f.Name) > duplicateWidth {
			duplicateWidth = len(f.Name)
		}
	}

	fmt.Printf("  %-*s  %-*s\n", originalWidth, originalHeader, duplicateWidth, duplicateHeader)
	for _, f := range deleted {
		fmt.Printf("  %-*s  %-*s\n", originalWidth, f.DuplicateOf, duplicateWidth, f.Name)
	}
}
