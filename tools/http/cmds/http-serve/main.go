package main

import "fmt"
import "net/http"
import "os"
import "strconv"

func showUsage() {

	fmt.Println("Usage: http-serve <port>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool serves the current folder via HTTP.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # Serve anything for other devices")
	fmt.Println("    cd ~/Videos;")
	fmt.Println("    http-serve 8080;")
	fmt.Println("")

}

func main() {

	var port int16

	if len(os.Args) == 2 {

		num, err := strconv.ParseInt(os.Args[1], 10, 16)

		if err == nil && num > 0 && num < 65535 {
			port = int16(num)
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a port")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}


	if port != 0 {

		root, err0 := os.Getwd()

		if err0 == nil {

			fs := os.DirFS(root)
			fsrv := http.FileServer(http.FS(fs))
			http.Handle("/", fsrv)

			fmt.Printf("Listening on http://localhost:%d\n", port)

			err1 := http.ListenAndServe(":" + strconv.Itoa(int(port)), nil)

			if err1 != nil {

				fmt.Fprintf(os.Stderr, "ERROR: Cannot listen on port %d\n", port)
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
				os.Exit(1)

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read current folder")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
