package main

import "net/http"
import "io"
import "os"
import "os/exec"
import "strings"

import "fmt"

func buildRequestURL(username string, password string, subdomain string, ipv6 string) string {

	// Optional parameter ttl=2000
	return fmt.Sprintf(
		"https://www.goip.de/setip?username=%s&password=%s&subdomain=%s&ip6=%s",
		username,
		password,
		subdomain,
		ipv6,
	)

}

func showUsage() {

	fmt.Println("Usage: dyndns-goip <username> <password> <subdomain.goip.de>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool sets the public IPv6 of this machine to an AAAA entry for a specified a goip.de subdomain.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Update with current public IPv6")
	fmt.Println("    dyndns-goip \"john_doe\" \"password123#special!\" \"myserver.goip.de\";")
	fmt.Println("")
}

func main() {

	var username string
	var password string
	var subdomain string

	if len(os.Args) == 4 {

		tmp1 := strings.TrimSpace(os.Args[1])

		if strings.HasPrefix(tmp1, "\"") && strings.HasSuffix(tmp1, "\"") {
			username = strings.TrimSpace(tmp1[1:len(tmp1)-1])
		} else if strings.HasPrefix(tmp1, "'") && strings.HasSuffix(tmp1, "'") {
			username = strings.TrimSpace(tmp1[1:len(tmp1)-1])
		} else if tmp1 != "" {
			username = strings.TrimSpace(tmp1)
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a username")

			showUsage()
			os.Exit(1)

		}

		tmp2 := strings.TrimSpace(os.Args[2])

		if strings.HasPrefix(tmp2, "\"") && strings.HasSuffix(tmp2, "\"") {
			password = strings.TrimSpace(tmp2[1:len(tmp2)-1])
		} else if strings.HasPrefix(tmp2, "'") && strings.HasSuffix(tmp2, "'") {
			password = strings.TrimSpace(tmp2[1:len(tmp2)-1])
		} else if tmp2 != "" {
			password = strings.TrimSpace(tmp2)
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Second argument has to be a password")

			showUsage()
			os.Exit(1)

		}

		tmp3 := strings.TrimSpace(os.Args[3])

		if strings.HasSuffix(tmp3, ".goip.de") {
			subdomain = tmp3
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Third argument has to be a subdomain")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if username != "" && password != "" && subdomain != "" {

		// ipv4s := make([]string, 0)
		// https://www.goip.de/setip?username=<username>&password=<password>&subdomain=<subdomain>&ip=<ip>&ttl=2000 

		cmd0 := exec.Command(
			"ip", "-6",
			"address", "show",
			"scope", "global",
		)
		buf0, err0 := cmd0.Output()

		if err0 == nil {

			ipv6  := ""
			lines := strings.Split(strings.TrimSpace(string(buf0)), "\n")

			for l := 0; l < len(lines); l++ {

				line := strings.TrimSpace(lines[l])

				if strings.HasPrefix(line, "inet6 ") && strings.Contains(line, " scope global") {

					tmp := strings.TrimSpace(line[6:strings.Index(line, " scope")])

					if tmp != "" && strings.HasSuffix(tmp, "/64") {

						tmp = tmp[0:strings.Index(tmp, "/")]
						ipv6 = tmp

						break

					}

				}

			}

			if ipv6 != "" {

				url := buildRequestURL(username, password, subdomain, ipv6)

				client := &http.Client{}

				request, err1 := http.NewRequest("GET", url, nil)

				if err1 == nil {

					fmt.Println("> Request URL: \"" + url + "\"")

					request.Header.Add("User-Agent", "Cookie's Forensics Tools (DynDNS Updater)")

					response, err2 := client.Do(request)

					if err2 == nil {

						if response.StatusCode == http.StatusOK {

							bytes, _ := io.ReadAll(response.Body)

							fmt.Println("> Updated AAAA entry of \"" + subdomain + "\" to IPv6 \"" + ipv6 + "\"")
							fmt.Println(string(bytes))
							os.Exit(0)

						} else {

							fmt.Fprintf(os.Stderr, "ERROR: HTTP response status: %d %s\n", response.StatusCode, response.Status)

							bytes, _ := io.ReadAll(response.Body)
							fmt.Fprintln(os.Stderr, string(bytes))

							os.Exit(1)

						}

					}

				} else {

					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
					os.Exit(1)

				}

			} else {

				fmt.Fprintln(os.Stderr, "ERROR: Cannot detect any global IPv6 for this machine")
				os.Exit(1)

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot execute iproute2 command")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
