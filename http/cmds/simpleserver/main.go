package main

import "fmt"
import "log"
import "net/http"
import "os"
import "strconv"

func showUsage() {
	fmt.Println("Usage: simplehttpserver [Port]")
	fmt.Println("")
	fmt.Println("The simplehttpserver will serve the current folder as assets via HTTP.")
}

func main() {

	var port int16 = 8080

	if len(os.Args) == 2 {

		num, err := strconv.ParseInt(os.Args[1], 10, 16)

		if err == nil {
			port = int16(num)
		}

	}

	root, err0 := os.Getwd()

	if err0 == nil {

		fs := os.DirFS(root)
		fsrv := http.FileServer(http.FS(fs))
		http.Handle("/", fsrv)

		fmt.Println("Listening on port " + strconv.Itoa(int(port)))

		err1 := http.ListenAndServe(":" + strconv.Itoa(int(port)), nil)

		if err1 != nil {
			log.Fatal(err1)
		}

	}

}
