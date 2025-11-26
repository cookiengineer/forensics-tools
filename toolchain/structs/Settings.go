package structs

import "os"
import "strings"

type Settings struct {

	Folder  string `json:"folder"`
	Debug   bool   `json:"debug"`

	// go build parameters
	Arch    string `json:"arch"` // GOARCH
	OS      string `json:"os"`   // GOOS

}

func ToSettings(folder string, debug bool) *Settings {

	var settings Settings

	settings.Folder = ""
	settings.Debug = false
	settings.Arch = "amd64"
	settings.OS = "linux"

	settings.SetDebug(debug)
	settings.SetFolder(folder)

	go_arch := os.Getenv("GOARCH")

	if go_arch != "" {
		settings.Arch = go_arch
	}

	go_os := os.Getenv("GOOS")

	if go_os != "" {
		settings.OS = go_os
	}

	return &settings

}

func (settings *Settings) SetDebug(value bool) {
	settings.Debug = value
}

func (settings *Settings) SetFolder(value string) bool {

	var result bool = false

	if strings.HasSuffix(value, "/toolchain") {
		value = value[0 : len(value)-10]
	} else if strings.HasSuffix(value, "/forensics-tools") {
		// Do Nothing
	} else {
		value = ""
	}

	if value != "" {

		settings.Folder = value
		result = true

	}

	return result

}
