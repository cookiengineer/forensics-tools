package main

import "forensics-toolchain/actions"
import "forensics-toolchain/structs"
import "os"
import "strings"

func showBuildUsage(console *structs.Console) {

	console.Info("")
	console.Info("toolchain/build")
	console.Info("")
	console.Group("Usage: build")
	console.Log("")
	console.GroupEnd("------")

}

func main() {

	console := structs.NewConsole(os.Stdout, os.Stderr, 0)

	folder, err := os.Getwd()

	if err == nil {

		if len(os.Args) == 2 {

			debug := strings.TrimSpace(strings.ToLower(os.Args[1]))

			if debug == "--debug" || debug == "--debug=true" {

				settings := structs.ToSettings(folder, true)
				result := actions.Build(console, settings)

				if result == true {

					console.Log("The installer has been built as ../build/install-forensics-tools")
					os.Exit(0)

				} else {
					os.Exit(1)
				}

			} else if debug == "--debug=false" {

				settings := structs.ToSettings(folder, false)
				result := actions.Build(console, settings)

				if result == true {

					console.Log("The installer has been built as ../build/install-forensics-tools")
					os.Exit(0)

				} else {
					os.Exit(1)
				}

			} else {

				showBuildUsage(console)
				os.Exit(1)

			}

		} else if len(os.Args) == 1 {

			settings := structs.ToSettings(folder, false)
			result := actions.Build(console, settings)

			if result == true {

				console.Log("The installer has been built as ../build/install-forensics-tools")
				os.Exit(0)

			} else {
				os.Exit(1)
			}

		} else {

			showBuildUsage(console)
			os.Exit(1)

		}

	} else {

		showBuildUsage(console)
		os.Exit(1)

	}

}
