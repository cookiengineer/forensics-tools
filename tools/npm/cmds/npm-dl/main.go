package main

import "npm/api"
import "fmt"
import "io"
import "net/http"
import "os"
import "regexp"
import "strings"

var expected_script_binaries []string = []string{
	"eslint",
	"jest",
	"npx",
	"prettier",
	"rimraf",
	"tsc",
	"yarn",
}

var dangerous_lifecycle_hooks []string = []string{
	"preinstall",
	"install",
	"postinstall",
	"preuninstall",
	"uninstall",
	"postuninstall",
	"prepublish",
	"publish",
	"postpublish",
}

func showUsage() {

	fmt.Println("Usage: npm-dl <package-name>")
	fmt.Println("Usage: npm-dl <package-name> <package-version>")
	fmt.Println("Usage: npm-dl <package-scope> <package-name> <package-version>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # download Shai Hulud worm samples")
	fmt.Println("    npm-dl web3-providers-http 4.1.0;")
	fmt.Println("    npm-dl @accordproject concerto-analysis 3.24.1;")
	fmt.Println("")

}

func main() {

	var package_scope string
	var package_name string
	var package_version string

	if len(os.Args) == 4 {

		pattern_scope := regexp.MustCompile(`^@([a-z0-9-_]+)$`)

		pattern_name := regexp.MustCompile(`^([a-z0-9-_]+)$`)

		pattern_semantic := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)
		pattern_prerelease := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)
		pattern_buildid := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)
		pattern_prerelease_buildid := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)

		tmp1 := strings.ToLower(os.Args[1])
		tmp2 := strings.ToLower(os.Args[2])
		tmp3 := strings.ToLower(os.Args[3])

		if pattern_scope.MatchString(tmp1) == true {
			package_scope = tmp1
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a package scope string")

			showUsage()
			os.Exit(1)

		}

		if pattern_name.MatchString(tmp2) == true {
			package_name = tmp2
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a package name string")

			showUsage()
			os.Exit(1)

		}

		switch {
		case pattern_semantic.MatchString(tmp3):
			package_version = tmp3
		case pattern_buildid.MatchString(tmp3):
			package_version = tmp3
		case pattern_prerelease.MatchString(tmp3):
			package_version = tmp3
		case pattern_prerelease_buildid.MatchString(tmp3):
			package_version = tmp3
		default:

			fmt.Fprintln(os.Stderr, "ERROR: Second argument has to be a semantic package version string")

			showUsage()
			os.Exit(1)

		}

	} else if len(os.Args) == 3 {

		pattern_name := regexp.MustCompile(`^([a-z0-9-_]+)$`)

		pattern_semantic := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)
		pattern_prerelease := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)
		pattern_buildid := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)
		pattern_prerelease_buildid := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)$`)

		tmp1 := strings.ToLower(os.Args[1])
		tmp2 := strings.ToLower(os.Args[2])

		if pattern_name.MatchString(tmp1) == true {
			package_name = tmp1
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a package name string")

			showUsage()
			os.Exit(1)

		}

		switch {
		case pattern_semantic.MatchString(tmp2):
			package_version = tmp2
		case pattern_buildid.MatchString(tmp2):
			package_version = tmp2
		case pattern_prerelease.MatchString(tmp2):
			package_version = tmp2
		case pattern_prerelease_buildid.MatchString(tmp2):
			package_version = tmp2
		default:

			fmt.Fprintln(os.Stderr, "ERROR: Second argument has to be a semantic package version string")

			showUsage()
			os.Exit(1)

		}

	} else if len(os.Args) == 2 {

		pattern_name := regexp.MustCompile(`^([a-z-_@]+)$`)

		tmp1 := strings.ToLower(os.Args[1])

		if pattern_name.MatchString(tmp1) == true {

			package_name = tmp1
			package_version = "latest"

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a package name string")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if package_name != "" && package_version != "" {

		schema, err0 := api.FetchPackage(package_scope, package_name, package_version)

		if err0 == nil {

			if len(schema.Scripts) > 0 {

				for name, cmdline := range schema.Scripts {

					found_hook := false
					found_binary := false

					for _, hook := range dangerous_lifecycle_hooks {

						if hook == name {
							found_hook = true
							break
						}

					}

					for _, binary := range expected_script_binaries {

						if strings.HasPrefix(cmdline, binary + " ") {
							found_binary = true
							break
						}

					}

					if found_binary == false {
						fmt.Printf("> Potentially dangerous hook \"%s\": \"%s\"\n", name, cmdline)
					} else if found_hook == true {
						fmt.Printf("> Potentially dangerous hook \"%s\": \"%s\"\n", name, cmdline)
					}

				}

			}

			if schema.Dist != nil && schema.Dist.Tarball != "" {

				client := &http.Client{}
				request, err1 := http.NewRequest("GET", schema.Dist.Tarball, nil)

				if err1 == nil {

					fmt.Printf("> Download URL: \"%s\"\n", schema.Dist.Tarball)

					request.Header.Set("User-Agent", "npm/10.5.0 node/v20.10.0 linux x64")
					request.Header.Set("Accept", "application/gzip")
					request.Header.Set("Accept-Encoding", "identity")

					response, err2 := client.Do(request)

					if err2 == nil {

						if response.StatusCode == http.StatusOK {

							output := fmt.Sprintf("%s_%s_%s.tgz", package_scope, package_name, package_version)
							file, err3 := os.Create(output)

							if err3 == nil {

								_, err4 := io.Copy(file, response.Body)

								if err4 == nil {

									fmt.Println("> Downloaded NPM package to \"" + output + "\"")
									os.Exit(0)

								} else {

									fmt.Fprintln(os.Stderr, "ERROR: Cannot write to \"" + output + "\"")
									os.Exit(1)

								}

							} else {

								fmt.Fprintln(os.Stderr, "ERROR: Cannot write to \"" + output + "\"")
								os.Exit(1)

							}

						} else {

							fmt.Fprintf(os.Stderr, "ERROR: HTTP response status: %d %s\n", response.StatusCode, response.Status)
							os.Exit(1)

						}

					} else {

						fmt.Fprintf(os.Stderr, "ERROR: %v\n", err2)
						os.Exit(1)

					}

				} else {

					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
					os.Exit(1)

				}

			}

		} else {
			fmt.Println("Error: ", err0.Error())
			os.Exit(1)
		}

	}

}
