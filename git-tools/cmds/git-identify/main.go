package main

import "fmt"
import "os"
import "strings"
import "git-tools/github"

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

func main() {

	var user string = ""

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "https://github.com/") {

			tmp := strings.TrimSpace(os.Args[1][19:])

			if strings.Contains(tmp, "/") {
				user = tmp[0:strings.Index(tmp, "/")]
			} else {
				user = strings.TrimSpace(tmp)
			}

		} else {
			user = os.Args[1]
		}

	} else {

		showHelp()
		os.Exit(1)

	}

	if user != "" {

		repos := github.ToRepositories(user)

		// TODO: Sort repositories by Datetime, earliest first
		// TODO: Filter repositories with IsFork == true
		// TODO: for each in repositories
		//       - ToCommits(api.Repository)
		//       - ParseCommit() into Name, EMail etc

		fmt.Println(user)
		fmt.Println(repos)

	}

}
