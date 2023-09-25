package github

import "git-tools/structs"
import "git-tools/github/api"
import "encoding/json"

func ToCommits(user string, repo string) []api.Commit {

	var scraper structs.Scraper

	scraper.Headers = map[string]string{
		"Accept": "application/json",
	}

	buffer := scraper.Request("https://api.github.com/repos/" + user + "/" + repo + "/commits?limit=100")

	var commits []api.Commit

	json.Unmarshal(buffer, &commits)

	return commits

}

