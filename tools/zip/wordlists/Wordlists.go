package wordlists

import "bytes"
import "compress/gzip"
import "embed"
import "fmt"
import "io"
import "io/fs"
import "strings"

//go:embed *
var embedded_fs embed.FS

var Wordlists map[string][]string

func init() {

	Wordlists = make(map[string][]string, 0)

	entries, err0 := fs.ReadDir(embedded_fs, ".")

	if err0 == nil {

		for _, entry := range entries {

			name := entry.Name()

			if strings.HasSuffix(name, ".txt.gz") {

				fmt.Println("> Reading embedded wordlist \"" + name + "\"")

				buffer, err1 := embedded_fs.ReadFile(name)

				if err1 == nil {

					buffer := bytes.NewBuffer(buffer)
					reader, err2 := gzip.NewReader(buffer)

					if err2 == nil {

						decompressed, err3 := io.ReadAll(reader)
						defer reader.Close()

						if err3 == nil {
							Wordlists[name] = strings.Split(strings.TrimSpace(string(decompressed)), "\n")
						}

					}

				}

			} else if strings.HasSuffix(name, ".txt") {

				fmt.Println("> Reading embedded wordlist \"" + name + "\"")

				buffer, err1 := embedded_fs.ReadFile(name)

				if err1 == nil {
					Wordlists[name] = strings.Split(strings.TrimSpace(string(buffer)), "\n")
				}

			}

		}

	}

}
