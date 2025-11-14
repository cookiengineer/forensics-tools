package main

import "dns/blocklists"
import "fmt"
import "os"
import "sort"
import "strings"

var Whitespace string = "                                                                        "
var Borderspace string = "------------------------------------------------------------------------"

func showUsage() {

	fmt.Println("Usage: dns-iscensored <pattern>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool searches known blocklists for a pattern to find out if the ISP is blocking it.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Search with suffix")
	fmt.Println("    dns-iscensored *.kinox.to; # exit code 1, blocked")
	fmt.Println("")
	fmt.Println("    # Search with midfix")
	fmt.Println("    dns-iscensored cookie*.engineer; # exit code 0, not blocked (hopefully)")
	fmt.Println("")
	fmt.Println("    # Search with prefix")
	fmt.Println("    dns-iscensored www4.kino*; # exit code 1, blocked")
	fmt.Println("")

}

func main() {

	var search_prefix string
	var search_midfix string
	var search_suffix string

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "*") {

			search_suffix = os.Args[1][1:]

		} else if strings.HasSuffix(os.Args[1], "*") {

			search_prefix = os.Args[1][0:len(os.Args[1])-1]

		} else if strings.Contains(os.Args[1], "*") {

			tmp := strings.Split(os.Args[1], "*")

			if len(tmp) == 2 {
				search_prefix = tmp[0]
				search_suffix = tmp[0]
			}

		} else {
			search_midfix = os.Args[1]
		}

	} else {

		showUsage()
		os.Exit(1)

	}


	if search_prefix != "" || search_midfix != "" || search_suffix != "" {

		filtered := make(map[string]bool)

		for _, blocklist := range blocklists.Blocklists {

			for domain, isblocked := range blocklist {

				if search_prefix != "" && search_suffix != "" {

					if strings.HasPrefix(domain, search_prefix) && strings.HasSuffix(domain, search_suffix) {
						filtered[domain] = isblocked
					}

				} else if search_prefix != "" {

					if strings.HasPrefix(domain, search_prefix) {
						filtered[domain] = isblocked
					}

				} else if search_suffix != "" {

					if strings.HasSuffix(domain, search_suffix) {
						filtered[domain] = isblocked
					}

				} else if search_midfix != "" {

					if strings.Contains(domain, search_midfix) {
						filtered[domain] = isblocked
					}

				}

			}

		}

		// Add entry for not blocked domain if it was not a pattern
		if strings.Contains(os.Args[1], "*") == false && len(filtered) == 0 {
			filtered[os.Args[1]] = false
		}

		if len(filtered) > 0 {

			domains := make([]string, 0)

			for domain, _ := range filtered {
				domains = append(domains, domain)
			}

			sort.Strings(domains)

			fmt.Println("Censored domains matching pattern \"" + os.Args[1] + "\":")
			fmt.Println("")

			maxlength := int(8)

			for domain, _ := range filtered {

				if len(domain) > maxlength {
					maxlength = len(domain)
				}

			}

			fmt.Println("| Domain" + Whitespace[0:(maxlength - 6)] + " | Blocked? |")
			fmt.Println("|-" + Borderspace[0:maxlength] + "-|----------|")

			found := false

			for _, domain := range domains {

				isblocked := filtered[domain]
				length := len(domain)
				offset := Whitespace[0:(maxlength - length)]

				if isblocked == true {
					fmt.Println("| " + domain + offset + " | true     |")
					found = true
				} else {
					fmt.Println("| " + domain + offset + " | false    |")
				}

			}

			fmt.Println("|-" + Borderspace[0:maxlength] + "-|----------|")

			if found == true {
				os.Exit(1)
			} else {
				os.Exit(0)
			}

		} else {

			fmt.Println("No domains matching \"" + os.Args[1] + "\" found")
			os.Exit(0)

		}

	}

}
