package main

import "torrent/structs"
import "torrent/trackers"
import "os"
import "strings"
import "fmt"

func showUsage() {

	fmt.Println("Usage: torrent-magnetify <magnet:url>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # output magnet:url with all default trackers")
	fmt.Println("    torrent-magnetify magnet:?xt=urn:btih:a55ac5f0580c53777cfa765d1f86e50dcac50fc0")
	fmt.Println("")

}

func main() {

	var magnet_link string

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "magnet:?xt=") {
			magnet_link = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a magnet:url")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if magnet_link != "" {

		magnet := structs.NewMagnet(os.Args[1])

		if magnet.ExactTopic != "" {

			for _, tracker := range trackers.Trackers {
				magnet.AddTracker(tracker)
			}

			fmt.Println(magnet.Render())
			os.Exit(0)

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Invalid magnet link")
			os.Exit(1)

		}

	}

}
