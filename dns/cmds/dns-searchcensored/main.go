package main

import _ "embed"
import "encoding/json"
import "fmt"
import "os"
import "sort"
import "strings"

//go:embed blocked_domains.json
var embedded_domains []byte

var Whitespace string = "                                                                        "
var Domains map[string]string

type Domain struct {
	AddedBy        string `json:"added_by"`
	Domain         string `json:"domain"`
	FirstBlockedOn string `json:"first_blocked_on"`
	// Site is always null?
}

func init() {

	Domains = make(map[string]string)

	var tmp []Domain

	err := json.Unmarshal(embedded_domains, &tmp)

	if err == nil {

		for t := 0; t < len(tmp); t++ {

			entry := tmp[t]

			Domains[entry.Domain] = entry.FirstBlockedOn

		}

	}

}

func main() {

	var filtered []string

	if len(os.Args) == 2 {

		var search_prefix string
		var search_midfix string
		var search_suffix string

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

		for domain := range Domains {

			if search_prefix != "" && search_suffix != "" {

				if strings.HasPrefix(domain, search_prefix) && strings.HasSuffix(domain, search_suffix) {
					filtered = append(filtered, domain)
				}

			} else if search_prefix != "" {

				if strings.HasPrefix(domain, search_prefix) {
					filtered = append(filtered, domain)
				}

			} else if search_suffix != "" {

				if strings.HasSuffix(domain, search_suffix) {
					filtered = append(filtered, domain)
				}

			} else if search_midfix != "" {

				if strings.Contains(domain, search_midfix) {
					filtered = append(filtered, domain)
				}

			}

		}

		sort.Strings(filtered)

		fmt.Println("Blocked Domains matching \"" + os.Args[1] + "\":")
		fmt.Println("")

	} else {

		for domain := range Domains {
			filtered = append(filtered, domain)
		}

		sort.Strings(filtered)

		fmt.Println("Blocked Domains matching \"*\":")
		fmt.Println("")

	}

	if len(filtered) > 0 {

		var maxlength int = 6

		for f := 0; f < len(filtered); f++ {

			if len(filtered[f]) > maxlength {
				maxlength = len(filtered[f])
			}

		}

		fmt.Println("Domain" + Whitespace[0:(maxlength - 6)] + " |    Date    | ")

		for f := 0; f < len(filtered); f++ {

			domain := filtered[f]
			timestamp := Domains[domain]

			length := len(domain)
			offset := Whitespace[0:(maxlength - length)]

			fmt.Println(domain + offset + " | " + timestamp + " | ")

		}

		os.Exit(0)

	} else {

		fmt.Println("(No Results)")
		os.Exit(1)

	}

}
