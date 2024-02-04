package github

import "git-tools/structs"
import "git-tools/github/api"
import "encoding/json"
import "strconv"

func ToRepositories(user string) []api.Repository {

	var scraper structs.Scraper

	scraper.Headers = map[string]string{
		"Accept": "application/json",
		"Token": "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",
		"User-Agent": "git-identify (Cookie Engineer's Forensics Tools)",
	}

	buffer := scraper.Request("https://api.github.com/users/" + user + "/repos?page=1")

	var repositories []api.Repository

	err := json.Unmarshal(buffer, &repositories)

	if err == nil && len(repositories) == 30 {

		for p := 2; p <= 10; p++ {

			var page = strconv.Itoa(p)
			var page_repositories []api.Repository

			page_buffer := scraper.Request("https://api.github.com/users/" + user + "/repos?page=" + page)

			if len(page_buffer) > 0 {

				err := json.Unmarshal(page_buffer, &page_repositories)

				if err == nil {

					for pr := 0; pr < len(page_repositories); pr++ {
						repositories = append(repositories, page_repositories[pr])
					}

					if len(page_repositories) < 30 {
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

	return repositories

}
