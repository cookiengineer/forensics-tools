package actions

import "forensics-toolchain/structs"
import "forensics-toolchain/utils"
import "bytes"
import "os/exec"
import "strings"
import "path/filepath"

func BuildTool(console *structs.Console, settings *structs.Settings, binary_name string, source_path string) bool {

	var result bool

	source_path = filepath.Clean(source_path)
	go_compiler := "go"
	ld_flags := ""

	if settings.Debug == false {
		go_compiler = "go"
		ld_flags = "-s -w"
	}

	if settings.Folder != "" && utils.IsProgram("env") && utils.IsProgram(go_compiler) {

		if strings.HasPrefix(source_path, "/") {
			source_path = strings.TrimPrefix(source_path, "/")
		}

		go_arch := settings.Arch
		go_os := settings.OS
		go_output := settings.Folder + "/toolchain/installer/programs/" + binary_name
		go_source := settings.Folder + "/" + source_path

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

		// /tools/whatever/cmds/whatever-cmd/main.go -> /tools/whatever
		tmp1 := strings.Split(source_path, string(filepath.Separator))
		tmp2 := filepath.Join(tmp1[0], tmp1[1])
		cmd1.Dir = settings.Folder + "/" + tmp2

		cmd1.Stdout = &stdout
		cmd1.Stderr = &stderr

		err1 := cmd1.Run()

		if err1 == nil {

			if settings.Debug == false {

				cmd2 := exec.Command(
					"strip",
					go_output,
				)
				cmd2.Dir = settings.Folder + "/toolchain/installer/programs"
				cmd2.Run()

			}

			console.Info("> " + binary_name + " / " + go_os + " / " + go_arch)

			result = true

		} else {

			console.Error("> " + binary_name + " / " + go_os + " / " + go_arch)

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

	return result

}
