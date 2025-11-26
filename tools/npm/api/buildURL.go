package api

import "fmt"
import "strings"

func buildURL(package_scope string, package_name string, package_version string) string {

	encoded_package_name := ""
	encoded_package_version := ""

	if strings.HasPrefix(package_scope, "@") {
		encoded_package_name = "%40" + strings.TrimPrefix(package_scope, "@") + "/" + package_name
	} else {
		encoded_package_name = package_name
	}

	if package_version != "" {
		encoded_package_version = package_version
	} else {
		encoded_package_version = "latest"
	}

	return fmt.Sprintf("https://registry.npmjs.org/%s/%s", encoded_package_name, encoded_package_version)

}
