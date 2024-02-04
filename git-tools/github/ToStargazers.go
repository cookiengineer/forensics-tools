package github

import "git-tools/structs"
import "git-tools/github/api"
import "encoding/json"
import "strconv"

func ToStargazers(user string, repo string) []api.User {

	scraper := structs.NewScraper(1)
	scraper.Headers = map[string]string{
		"Accept": "application/json",
		// "Accept-Encoding": "identity",
		// "Accept-Language": "en-US,en;q=0.5",
		// "Cache-Control": "no-cache",
		// "Pragma": "no-cache",
		// "Sec-Fetch-Dest": "document",
		// "Sec-Fetch-Mode": "navigate",
		// "Sec-Fetch-Site": "cross-site",
		// "Sec-Fetch-User": "?1",
		// "Upgrade-Insecure-Requests": "1",
		// "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/113.0",
		"Token": "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",
		"User-Agent": "git-identify (Cookie Engineer's Forensics Tools)",
	}

	buffer := scraper.Request("https://api.github.com/repos/" + user + "/" + repo + "/stargazers?page=1")

	var users []api.User

	err := json.Unmarshal(buffer, &users)

	if err == nil && len(users) == 30 {

		for p := 2; p <= 50; p++ {

			var page = strconv.Itoa(p)
			var page_users []api.User

			page_buffer := scraper.Request("https://api.github.com/repos/" + user + "/" + repo + "/stargazers?page=" + page)

			if len(page_buffer) > 0 {

				err := json.Unmarshal(page_buffer, &page_users)

				if err == nil {

					for pu := 0; pu < len(page_users); pu++ {
						users = append(users, page_users[pu])
					}

					if len(page_users) < 30 {
						break
					}

				} else {
					break
				}

			} else {
				break
			}

		}

	}

	return users

}

