package main

import "encoding/json"
import "fmt"
import "os"
import "path"
import "strings"
import "git-tools/github"
import "git-tools/github/api"

func showHelp() {

	fmt.Println("git-identify <github-username>")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("    Identify a specific GitHub username with their aliases and email addresses.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    git-identify name-of-doxxer > name-of-doxxer.json;")
	fmt.Println("")

}

func readUser(name string, file string) api.User {

	var user = api.NewUser(name)

	stat, err1 := os.Stat(file)

	if err1 == nil && stat.IsDir() == false {

		buffer, err2 := os.ReadFile(file)

		if err2 == nil && len(buffer) > 2 {
			json.Unmarshal(buffer, &user)
		}

	}

	return user

}

func writeFile(file string, buffer []byte) bool {

	var result bool = false

	folder := path.Dir(file)
	stat1, err1 := os.Stat(folder)

	if err1 == nil && stat1.IsDir() {

		err2 := os.WriteFile(file, buffer, 0666)

		if err2 == nil {
			result = true
		}

	} else {

		err2 := os.MkdirAll(folder, 0755)

		if err2 == nil {

			err3 := os.WriteFile(file, buffer, 0666)

			if err3 == nil {
				result = true
			}

		}

	}

	return result

}

func main() {

	var user string = ""
	var repo string = ""

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "https://github.com/") {

			tmp := strings.TrimSpace(os.Args[1][19:])

			if strings.Contains(tmp, "/") {

				user = strings.TrimSpace(tmp[0:strings.Index(tmp, "/")])
				repo = strings.TrimSpace(tmp[strings.Index(tmp, "/")+1:])

				if strings.Contains(repo, "/") {
					repo = strings.TrimSpace(repo[0:strings.Index(repo, "/")])
				}

			} else {
				user = strings.TrimSpace(tmp)
			}

		} else {
			user = strings.TrimSpace(os.Args[1])
		}
	
	} else {

		showHelp()
		os.Exit(1)

	}

	cwd, err := os.Getwd()

	if err == nil {

		if user != "" && repo != "" {

			fmt.Println("TODO: watchers and contributors")

			stargazers := github.ToStargazers(user, repo)

			buffer, _ := json.MarshalIndent(stargazers, "", "\t")
			writeFile(cwd + "/" + user + "_" + repo + "_stargazers.json", buffer)

			var usernames []string

			for s := 0; s < len(stargazers); s++ {
				usernames = append(usernames, stargazers[s].Name)
			}

			writeFile(cwd + "/" + user + "_" + repo + "_stargazers.txt", []byte(strings.Join(usernames, "\n")))

		} else if user != "" {

			var user = readUser(user, cwd + "/" + user + ".json")

			repositories := github.ToRepositories(user.Name)

			for r := 0; r < len(repositories); r++ {

				repository := repositories[r]

				if repository.IsFork == false && user.HasRepository(repository.Name) == false {

					user.AddRepository(repository.Name)

					commits := github.ToCommits(user.Name, repository.Name)

					for c := 0; c < len(commits); c++ {

						commit := commits[c]

						if commit.Commit.Author.Name != "" && commit.Commit.Author.Email != "" {

							if strings.HasSuffix(commit.Commit.Author.Email, "@users.noreply.github.com") == false {
								user.AddAlias(commit.Commit.Author.Name)
								user.AddEmail(commit.Commit.Author.Email)
							}

						}

						if commit.Commit.Committer.Name != "" && commit.Commit.Committer.Email != "" {

							if strings.HasSuffix(commit.Commit.Committer.Email, "@users.noreply.github.com") == false {
								user.AddAlias(commit.Commit.Committer.Name)
								user.AddEmail(commit.Commit.Committer.Email)
							}

						}

					}

					buffer, _ := json.MarshalIndent(user, "", "\t")
					writeFile(cwd + "/" + user.Name + ".json", buffer)

				}

			}

			buffer, _ := json.MarshalIndent(user, "", "\t")
			writeFile(cwd + "/" + user.Name + ".json", buffer)

		}

	}

}
