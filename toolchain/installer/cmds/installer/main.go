package main

import "installer"
import io_fs "io/fs"
import "fmt"
import "os"
import "path/filepath"
import "strings"

func main() {

	prefix := "/usr/local"

	if os.Geteuid() != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: Please run this program as root.\n")
		os.Exit(1)
	}

	tmp1 := os.Getenv("PREFIX")

	if strings.HasPrefix(tmp1, "/") {
		prefix = tmp1
	}

	bin_folder := filepath.Join(prefix, "bin")

	if err := os.MkdirAll(bin_folder, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create directory %s: %v", bin_folder, err)
		os.Exit(1)
	}

	entries, err0 := io_fs.ReadDir(installer.Filesystem, "programs")

	if err0 == nil {

		has_errors := false

		for _, entry := range entries {

			filename := entry.Name()

			if strings.HasPrefix(filename, ".") == false && entry.IsDir() == false {

				src_path := filepath.Join("programs", filename)
				dest_path := filepath.Join(bin_folder, filename)

				bytes, err1 := installer.Filesystem.ReadFile(src_path)

				if err1 == nil {

					err2 := os.WriteFile(dest_path, bytes, 0755)

					if err2 == nil {
						fmt.Fprintf(os.Stdout, "> Installed %s as %s\n", filename, dest_path)
					} else {
						fmt.Fprintf(os.Stderr, "ERROR: Failed to install program %s: %v\n", filename, err2)
						has_errors = true
					}

				}

			}

		}

		if has_errors == true {
			os.Exit(1)
		} else {
			os.Exit(0)
		}

	} else {
		os.Exit(1)
	}

}
