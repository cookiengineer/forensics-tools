package main

import "fmt"
import "io"
import "net/http"
import "net/url"
import "os"
import "regexp"
import "runtime"
import "strings"

func buildDownloadURL(extension_id string) string {

	update_url := "https://clients2.google.com/service/update2/crx?%s"

	params := url.Values{}
	params.Set("response", "redirect")
	params.Set("acceptformat", "crx2,crx3")
	params.Set("prodversion", "141.0") // Chrome Version needs to be updated

	os := runtime.GOOS
	arch := ""
	nacl_arch := ""

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86"
		nacl_arch = "x86-64"
	case "386":
		arch = "x86"
		nacl_arch = "x86-32"
	case "arm64":
		arch = "arm"
		nacl_arch = "arm64"
	case "arm":
		arch = "arm"
		nacl_arch = "arm"
	default:
		arch = "x64"
		nacl_arch = "x86-64"
	}

	params.Set("os", os)
	params.Set("arch", arch)
	params.Set("nacl_arch", nacl_arch)
	params.Set("prod", "chrome")
	params.Set("prodchannel", "unknown")
	params.Set("lang", "en-US")

	params.Set("x", fmt.Sprintf("id=%s&installsource=ondemand&uc", extension_id))

	return fmt.Sprintf(update_url, params.Encode())

}

func showUsage() {

	fmt.Println("Usage: crx-dl <https://chromewebstore-url>")
	fmt.Println("Usage: crx-dl <extension-id>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # download ublock origin lite")
	fmt.Println("    crx-dl https://chromewebstore.google.com/detail/ublock-origin-lite/ddkjiahejlhfcafbddmgiahcphecmpfh")
	fmt.Println("    crx-dl ddkjiahejlhfcafbddmgiahcphecmpfh")
	fmt.Println("")

}

func main() {

	var extension_id string

	if len(os.Args) == 2 {

		pattern := regexp.MustCompile(`^[a-z]{32}$`)

		if strings.HasPrefix(os.Args[1], "https://chromewebstore.google.com/detail/") {

			url := strings.TrimSpace(os.Args[1])

			if strings.Contains(url, "?") {

				tmp := strings.TrimSpace(url[strings.LastIndex(url, "/")+1:strings.Index(url, "?")])

				if pattern.MatchString(tmp) == true {
					extension_id = tmp
				}

			} else {

				tmp := strings.TrimSpace(url[strings.LastIndex(url, "/")+1:])


				if pattern.MatchString(tmp) == true {
					extension_id = tmp
				}

			}

		} else if pattern.MatchString(os.Args[1]) == true {
			extension_id = os.Args[1]
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a CRX url or CRX identifier")

			showUsage()
			os.Exit(1)

		}

	}

	if extension_id != "" {

		fmt.Println("> Extension ID: \"" + extension_id + "\"")

		url := buildDownloadURL(extension_id)

		client := &http.Client{}

		request, err1 := http.NewRequest("GET", url, nil)

		if err1 == nil {

			fmt.Println("> Download URL: \"" + url + "\"")

			request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")
			request.Header.Set("Accept", "*/*")
			request.Header.Set("Connection", "keep-alive")

			response, err2 := client.Do(request)

			if err2 == nil {

				if response.StatusCode == http.StatusOK {

					output := extension_id + ".crx"
					file, err3 := os.Create(output)

					if err3 == nil {

						_, err4 := io.Copy(file, response.Body)

						if err4 == nil {

							fmt.Println("> Downloaded CRX file to \"" + output + "\"")
							os.Exit(0)

						} else {

							fmt.Fprintln(os.Stderr, "ERROR: Cannot write to \"" + output + "\"")
							os.Exit(1)

						}

					} else {

						fmt.Fprintln(os.Stderr, "ERROR: Cannot write to \"" + output + "\"")
						os.Exit(1)

					}

				} else {

					fmt.Fprintf(os.Stderr, "ERROR: HTTP response status: %d %s\n", response.StatusCode, response.Status)
					os.Exit(1)

				}

			} else {

				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err2)
				os.Exit(1)

			}

		} else {

			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
			os.Exit(1)

		}

	} else {

		fmt.Fprintln(os.Stderr, "ERROR: Unsupported URL or CRX id format")

		showUsage()
		os.Exit(1)

	}

}

