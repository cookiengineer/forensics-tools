package trackers

import _ "embed"
import "slices"
import "strings"

var Trackers []string

//go:embed Trackers.txt
var embedded_Trackers []byte

func init() {

	lines := strings.Split(string(embedded_Trackers), "\n")

	for l := 0; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])

		if strings.HasPrefix(line, "udp://") || strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {

			if !slices.Contains(Trackers, line) {
				Trackers = append(Trackers, line)
			}

		}

	}

}
