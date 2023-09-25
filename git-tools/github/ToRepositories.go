package github

import "git-tools/structs"
import "git-tools/github/api"
import "encoding/json"

func ToRepositories(user string) []api.Repository {

	var scraper structs.Scraper

	scraper.Headers = map[string]string{
		"Accept": "application/json",
	}

	buffer := scraper.Request("https://api.github.com/users/" + user + "/repos?limit=100")

	var repositories []api.Repository

	json.Unmarshal(buffer, &repositories)

	return repositories

}
