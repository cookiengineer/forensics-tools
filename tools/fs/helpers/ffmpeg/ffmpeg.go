package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ffmpegBin and ffprobeBin are resolved once; if unavailable the package's
// functions return a descriptive error.
var (
	ffmpegBin  = lookup("ffmpeg")
	ffprobeBin = lookup("ffprobe")
)

func lookup(name string) string {
	bin, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return bin
}

// Available reports whether both ffmpeg and ffprobe are present.
func Available() bool {
	return ffmpegBin != "" && ffprobeBin != ""
}

// errUnavailable is returned when the required binary is not on PATH.
var errUnavailable = errors.New("ffmpeg/ffprobe not found on PATH")

// ExtractFrame decodes a single frame at seekSeconds and returns it encoded as
// PNG. The bytes can be fed to an image decoder (e.g. a perceptual hasher).
func ExtractFrame(path string, seekSeconds float64) ([]byte, error) {
	if ffmpegBin == "" {
		return nil, errUnavailable
	}
	args := []string{
		"-v", "error",
		"-ss", strconv.FormatFloat(seekSeconds, 'f', 3, 64),
		"-i", path,
		"-frames:v", "1",
		"-f", "image2pipe",
		"-c:v", "png",
		"-",
	}
	cmd := exec.Command(ffmpegBin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg frame extract %s: %w", path, err)
	}
	if out.Len() == 0 {
		return nil, fmt.Errorf("ffmpeg frame extract %s: empty output", path)
	}
	return out.Bytes(), nil
}

// ProbeDuration returns the duration of a video file in seconds.
func ProbeDuration(path string) (float64, error) {
	if ffprobeBin == "" {
		return 0, errUnavailable
	}
	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		path,
	}
	out, err := exec.Command(ffprobeBin, args...).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe duration %s: %w", path, err)
	}
	s := strings.TrimSpace(string(out))
	if s == "" || s == "N/A" {
		return 0, fmt.Errorf("ffprobe duration %s: no duration", path)
	}
	d, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("ffprobe duration %s: %w", path, err)
	}
	return d, nil
}
