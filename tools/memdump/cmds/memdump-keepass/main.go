package main

import "fmt"
import "os"
import "strings"

func showUsage() {

	fmt.Println("Usage: memdump-keepass <memdump.mem>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool searches for a KeePass password in a raw memdump.")
	fmt.Println("    Supported formats are raw memory dumps (e.g. WinDbg, MiniDumpWriteDump, dd)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Find KeePass password in memory dump")
	fmt.Println("    memdump-keepass ~/Downloads/evidence-1337.mem")
	fmt.Println("")

}
func main() {

	var memdump_file string

	if len(os.Args) == 2 {

		file := strings.TrimSpace(os.Args[1])

		if strings.HasSuffix(file, ".dmp") {
			memdump_file = file
		} else if strings.HasSuffix(file, ".raw") {
			memdump_file = file
		} else if strings.HasSuffix(file, ".mem") {
			memdump_file = file
		} else if strings.HasSuffix(file, ".bin") {
			memdump_file = file
		} else if strings.HasSuffix(file, ".mdmp") || strings.HasSuffix(file, ".memdump") {
			memdump_file = file
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a memory dump file")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if memdump_file != "" {

		memdump_buffer, err0 := os.ReadFile(memdump_file)

		if err0 == nil {

			candidates := make([]int, 0)

			for b := 0; b < len(memdump_buffer); b++ {

				var offset = b

				if memdump_buffer[offset] == 0xCF && memdump_buffer[offset + 1] == 0x25 {
					candidates = append(candidates, offset)
				}

			}

			if len(candidates) > 0 {

				for _, candidate := range candidates {

					offset := candidate + 2

					fmt.Println("> Found password candidate at offset 0x%x\n", offset)

					// assume less than 64 bytes password length
					for length := 0; length < 64; length++ {

						// 0x20 - 0xFF are valid characters
						if memdump_buffer[offset + length] >= 0x20 && memdump_buffer[offset + length] <= 0xFF {
							fmt.Printf("> Password might be \"%s\" with length %d\n", string(memdump_buffer[offset:offset + length]), length)
						} else {
							break
						}

					}

				}

				os.Exit(0)

			} else {

				fmt.Println("> No password candidates found")
				os.Exit(0)

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read file \"" + memdump_file + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
