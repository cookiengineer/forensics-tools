package main

import "torrent-tools/structs"
import "torrent-tools/trackers"
import "os"
import "strings"
import "fmt"

func main() {

	var magnet structs.Magnet

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "magnet:?xt=") {

			magnet = structs.NewMagnet(os.Args[1])

			if magnet.ExactTopic != "" {

				for t := 0; t < len(trackers.Trackers); t++ {
					magnet.AddTracker(trackers.Trackers[t])
				}

				fmt.Println("Use this link here:")
				fmt.Println("")
				fmt.Println(magnet.Render())

				os.Exit(0)

			} else {
				fmt.Println("Error: Invalid magnet link")
				os.Exit(1)
			}

		}

	}

}
