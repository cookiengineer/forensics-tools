package main

import "net/http"
import "io"
import "os"
import "os/exec"
import "strings"

import "fmt"

func showUsage() {
	fmt.Println("goip-updater usage:")
	fmt.Println("")
	fmt.Println("    --username=\"...\"   Account Username")
	fmt.Println("    --password=\"...\"   Account Password")
	fmt.Println("    --subdomain=\"...\"  GoIP Subdomain (with .goip.de)")
	fmt.Println("")
	fmt.Println("example:")
	fmt.Println("")
	fmt.Println("    goip-updater --username=\"john_doe\" --password=\"password123\" --subdomain=\"myserver.goip.de\";")
	fmt.Println("")
}

func main() {

	username := ""
	password := ""
	subdomain := ""

	if len(os.Args) > 0 {

		for o := 1; o < len(os.Args); o++ {

			tmp1 := os.Args[o]

			if strings.HasPrefix(tmp1, "--username=") {

				tmp2 := strings.TrimSpace(tmp1[11:])

				if strings.HasPrefix(tmp2, "\"") && strings.HasSuffix(tmp2, "\"") {
					tmp2 = strings.TrimSpace(tmp2[1:len(tmp2)-1])
				}

				if tmp2 != "" {
					username = tmp2
				}

			} else if strings.HasPrefix(tmp1, "--password=") {

				tmp2 := strings.TrimSpace(tmp1[11:])

				if strings.HasPrefix(tmp2, "\"") && strings.HasSuffix(tmp2, "\"") {
					tmp2 = tmp2[1:len(tmp2)-1]
				}

				if tmp2 != "" {
					password = tmp2
				}

			} else if strings.HasPrefix(tmp1, "--subdomain=") {

				tmp2 := strings.TrimSpace(tmp1[12:])

				if strings.HasPrefix(tmp2, "\"") && strings.HasSuffix(tmp2, "\"") {
					tmp2 = tmp2[1:len(tmp2)-1]
				}

				if tmp2 != "" && strings.HasSuffix(tmp2, ".goip.de") {
					subdomain = tmp2
				}

			}

		}

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

				client := &http.Client{}

				request, err1 := http.NewRequest("GET", strings.Join([]string{
					"https://www.goip.de/setip",
					"?username=" + username,
					"&password=" + password,
					"&subdomain=" + subdomain,
					"&ip6=" + ipv6,
					// "&ttl=2000",
				}, ""), nil)

				if err1 == nil {

					request.Header.Add("User-Agent", "Cookie's DynDNS Updater")

					response, err2 := client.Do(request)

					if err2 == nil {

						bytes, _ := io.ReadAll(response.Body)
						fmt.Println(string(bytes))

						fmt.Println("Done!")
						os.Exit(0)

					} else {

						if response.StatusCode != http.StatusOK {

							bytes, _ := io.ReadAll(response.Body)
							fmt.Println(string(bytes))

						}

						os.Exit(1)

					}

				}

			} else {
				fmt.Println("Error: Could not detect any global IPv6 for this machine")
				os.Exit(1)
			}

		} else {
			fmt.Println("Error: No iproute2 package installed?")
			os.Exit(1)
		}

	} else {
		showUsage()
		os.Exit(1)
	}

}
