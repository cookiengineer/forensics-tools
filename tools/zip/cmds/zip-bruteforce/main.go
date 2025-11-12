package main

import "zip/wordlists"
import "github.com/yeka/zip"
import _ "embed"
import "fmt"
import "io"
import "os"
import "strings"
import "sync"

func showUsage() {

	fmt.Println("Usage: zip-bruteforce <file.zip>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Show password if any password of embedded wordlist matches")
	fmt.Println("    zip-bruteforce ~/Downloads/evidence-1337.zip")
	fmt.Println("")

}

func decryptFile(reader *zip.ReadCloser, password string) string {

	var result string

	for _, file := range reader.File {

		if file.IsEncrypted() {
			file.SetPassword(password)
		}

		stream, err1 := file.Open()

		if err1 == nil {

			_, err2 := io.ReadAll(stream)

			if err2 == nil {
				result = password
			}

			defer stream.Close()

		}

	}

	return result

}

func main() {

	var zip_file string

	if len(os.Args) == 2 {

		if strings.HasSuffix(os.Args[1], ".zip") {
			zip_file = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a ZIP file")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if zip_file != "" {

		channel := make(chan string)

		reader, err1 := zip.OpenReader(zip_file)

		if err1 == nil {

			waitgroup := sync.WaitGroup{}

			for wordlist, passwords := range wordlists.Wordlists {

				fmt.Println("")
				fmt.Printf("> Wordlist \"%s\"", wordlist)
				fmt.Println("")

				for p, password := range passwords {

					waitgroup.Add(1)

					go func (password string) {

						defer waitgroup.Done()
						result := decryptFile(reader, password)
						channel <- result

					}(password)

					result := <-channel

					if result != "" {

						fmt.Printf("\r\033[K> Tried out password \"%s\" (%d of %d)", password, p+1, len(passwords))
						fmt.Println("")
						fmt.Printf("> Password is \"%s\"\n", password)

						os.Exit(0)

						break

					} else {
						fmt.Printf("\r\033[K> Tried out password \"%s\" (%d of %d)", password, p+1, len(passwords))
					}

					waitgroup.Wait()

				}

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read file \"" + zip_file + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
			os.Exit(1)

		}

	}

}
