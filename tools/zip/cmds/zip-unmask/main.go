package main

import "fmt"
import "os"
import "strings"

var zip_headers [][]byte = [][]byte{
	// See: https://en.wikipedia.org/wiki/ZIP_(file_format)
	[]byte("PK\x03\x04"),
	[]byte("PK\x05\x06"),
	[]byte("PK\x07\x08"),
}

func showUsage() {

	fmt.Println("Usage: zip-unmask <xor-masked-file.zip.crypt>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool will try all possible XOR bitmasks and create unmasked .zip files in the same folder.")
	fmt.Println("    Supported are only XOR bitmask keys lower than or equal 32 bits.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Find the masking key and decrypt ZIP file as ~/Downloads/evidence-1337.zip")
	fmt.Println("    zip-unmask ~/Downloads/evidence-1337.zip.crypt")
	fmt.Println("")

}


func main() {

	var encrypted_file string

	if len(os.Args) == 2 {

		if strings.HasSuffix(os.Args[1], ".zip.crypt") {
			encrypted_file = strings.TrimSpace(os.Args[1])
		} else if strings.HasSuffix(os.Args[1], ".zip") {
			encrypted_file = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be an XOR encrypted ZIP file")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if encrypted_file != "" {

		encrypted_buffer, err0 := os.ReadFile(encrypted_file)

		if err0 == nil {

			encrypted_header := encrypted_buffer[0:4]
			possible_keys := make([][]byte, 0)
			possible_buffers := make([][]byte, 0)

			fmt.Printf("> Encrypted ZIP header: \"%x\"", encrypted_header)

			fmt.Println("> Generating possible decryption keys...")

			for _, expected_header := range zip_headers {

				possible_key := make([]byte, 0)

				for e := 0; e < len(expected_header); e++ {
					possible_key = append(possible_key, expected_header[e] ^ encrypted_header[e])
				}

				possible_keys = append(possible_keys, possible_key)

			}

			for _, possible_key := range possible_keys {

				possible_buffer := make([]byte, 0)

				for b := 0; b < len(encrypted_buffer); b++ {
					possible_buffer = append(possible_buffer, encrypted_buffer[b] ^ possible_key[b % 4])
				}

				possible_buffers = append(possible_buffers, possible_buffer)

			}

			for p, buffer := range possible_buffers {

				filename := fmt.Sprintf("unmasked-%d.zip", p)
				err := os.WriteFile(filename, buffer, 0666)

				if err == nil {

					fmt.Println("> Possible decrypted ZIP file unmasked as \"" + filename + "\"")

				} else {

					fmt.Fprintln(os.Stderr, "ERROR: Cannot write file \"" + filename + "\"")
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)

				}

			}

			os.Exit(0)

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read file \"" + encrypted_file + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
