package github

import "git-tools/structs"
import "git-tools/github/api"
import "encoding/json"
import "strconv"

func ToCommits(user string, repo string) []api.Commit {

	scraper := structs.NewScraper(1)
	scraper.Headers = map[string]string{
		"Accept": "application/json",
		"Token": "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",
		"User-Agent": "git-identify (Cookie Engineer's Forensics Tools)",
	}

	buffer := scraper.Request("https://api.github.com/repos/" + user + "/" + repo + "/commits?page=1")

	var commits []api.Commit

	err := json.Unmarshal(buffer, &commits)

	if err == nil && len(commits) == 30 {

		for p := 2; p <= 10; p++ {

			var page = strconv.Itoa(p)
			var page_commits []api.Commit

			page_buffer := scraper.Request("https://api.github.com/repos/" + user + "/" + repo + "/commits?page=" + page)

			if len(page_buffer) > 0 {

				err := json.Unmarshal(page_buffer, &page_commits)

				if err == nil {

					for pc := 0; pc < len(page_commits); pc++ {
						commits = append(commits, page_commits[pc])
					}

					if len(page_commits) < 30 {
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

	return commits

}

