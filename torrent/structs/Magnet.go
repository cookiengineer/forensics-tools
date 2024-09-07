package structs

import "torrent/utils"
import "slices"
import "strings"
import "fmt"

type Magnet struct {
	ExactTopic  string   `url:"xt"`
	ExactSource string   `url:"xs"`
	ExactLength string   `url:"xl"`
	DisplayName string   `url:"dn"`
	Trackers    []string `url:"tr"`
	WebSeeds    []string `url:"ws"`
}

func NewMagnet(raw string) Magnet {

	var magnet Magnet

	magnet.Trackers = make([]string, 0)
	magnet.WebSeeds = make([]string, 0)

	magnet.Parse(raw)

	return magnet

}

func (magnet *Magnet) Parse(raw string) {

	if strings.HasPrefix(raw, "magnet:?") {

		tmp1 := strings.TrimSpace(raw[8:])
		parameters := strings.Split(tmp1, "&")

		for p := 0; p < len(parameters); p++ {

			parameter := parameters[p]
			tmp2 := strings.Split(parameter, "=")

			if len(tmp2) == 2 {

				if tmp2[0] == "xt" {

					fmt.Println(tmp2[0], tmp2[1])

					if strings.HasPrefix(tmp2[1], "urn:btih:") {
						magnet.ExactTopic = tmp2[1]
					}

				} else if tmp2[0] == "xl" {

					if utils.IsNumber(tmp2[1]) {
						magnet.ExactLength = tmp2[1]
					}

				} else if tmp2[0] == "xs" {

					if strings.HasPrefix(tmp2[1], "http://") || strings.HasPrefix(tmp2[1], "https://") {
						magnet.ExactSource = tmp2[1]
					}

				} else if tmp2[0] == "dn" {

					name := utils.ToFilename(utils.ToASCII(tmp2[1]))

					if name != "" {
						magnet.DisplayName = name
					}

				} else if tmp2[0] == "tr" {

					if strings.HasPrefix(tmp2[1], "http://") || strings.HasPrefix(tmp2[1], "https://") || strings.HasPrefix(tmp2[1], "udp://") {

						if !slices.Contains(magnet.Trackers, tmp2[1]) {
							magnet.Trackers = append(magnet.Trackers, tmp2[1])
						}

					}

				} else if tmp2[0] == "ws" {

					if strings.HasPrefix(tmp2[1], "http://") || strings.HasPrefix(tmp2[1], "https://") {

						if !slices.Contains(magnet.Trackers, tmp2[1]) {
							magnet.WebSeeds = append(magnet.WebSeeds, tmp2[1])
						}

					}

				}

			}

		}

	}

}

func (magnet *Magnet) AddTracker(tracker string) {

	if strings.HasPrefix(tracker, "http://") || strings.HasPrefix(tracker, "https://") || strings.HasPrefix(tracker, "udp://") {

		if !slices.Contains(magnet.Trackers, tracker) {
			magnet.Trackers = append(magnet.Trackers, tracker)
		}

	}

}

func (magnet *Magnet) Render() string {

	var result string

	if strings.HasPrefix(magnet.ExactTopic, "urn:btih:") {

		result = "magnet:?xt=" + magnet.ExactTopic

		if magnet.DisplayName != "" {
			result += "&dn=" + magnet.DisplayName
		}

		// TODO: ExactSource
		// TODO: ExactLength

		if len(magnet.Trackers) > 0 {

			for t := 0; t < len(magnet.Trackers); t++ {
				result += "&tr=" + magnet.Trackers[t]
			}

		}

		if len(magnet.WebSeeds) > 0 {

			for w := 0; w < len(magnet.WebSeeds); w++ {
				result += "&ws=" + magnet.WebSeeds[w]
			}

		}

	}

	return result

}
