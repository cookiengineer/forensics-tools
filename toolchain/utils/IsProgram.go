package utils

import "os"

func IsProgram(name string) bool {

	var result bool = false

	folders := []string{
		"/sbin",
		"/usr/bin",
		"/usr/local/bin",
	}

	for f := 0; f < len(folders); f++ {

		folder := folders[f]

		stat, err := os.Stat(folder + "/" + name)

		if err == nil && stat.IsDir() == false {
			result = true
			break
		}

	}

	return result

}
