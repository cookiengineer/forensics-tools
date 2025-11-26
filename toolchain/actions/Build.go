package actions

import "forensics-toolchain/structs"
import "forensics-toolchain/utils"
import "bytes"
import "os/exec"
import "strings"

func Build(console *structs.Console, settings *structs.Settings) bool {

	var result bool

	console.Group("Build")

	go_compiler := "go"
	ld_flags := ""

	if settings.Debug == false {
		go_compiler = "go"
		ld_flags = "-s -w"
	}

	if settings.Folder != "" && utils.IsProgram("env") && utils.IsProgram(go_compiler) {

		for name, path := range utils.Tools {
			BuildTool(console, settings, name, path)
		}

		go_arch := settings.Arch
		go_os := settings.OS
		go_output := settings.Folder + "/build/install-forensics-tools"
		go_source := settings.Folder + "/toolchain/installer/cmds/installer/main.go"

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd1 := exec.Command(
			"env",
			"CGO_ENABLED=0",
			"GOOS=" + go_os,
			"GOARCH=" + go_arch,
			go_compiler,
			"build",
			"-ldflags",
			ld_flags,
			"-o",
			go_output,
			go_source,
		)
		cmd1.Dir = settings.Folder + "/toolchain/installer"

		cmd1.Stdout = &stdout
		cmd1.Stderr = &stderr

		err1 := cmd1.Run()

		if err1 == nil {

			if settings.Debug == false {

				cmd2 := exec.Command(
					"strip",
					go_output,
				)
				cmd2.Dir = settings.Folder + "/build"
				cmd2.Run()

			}

			console.Info("> installer / " + go_os + " / " + go_arch)

			result = true

		} else {

			console.Error("> installer / " + go_os + " / " + go_arch)

			stdout_message := strings.TrimSpace(string(stdout.Bytes()))
			stderr_message := strings.TrimSpace(string(stderr.Bytes()))

			if stdout_message != "" {
				console.Error(stdout_message)
			}

			if stderr_message != "" {
				console.Error(stderr_message)
			}

			result = false

		}

	} else {

		if utils.IsProgram("env") == false {
			console.Error("Cannot find \"env\": Please execute \"sudo pacman -S coreutils\"")
		}

		if utils.IsProgram("go") == false {
			console.Error("Cannot find \"go\": Please execute \"sudo pacman -S go\"")
		}

	}

	console.GroupEnd("Build")

	return result

}
