package main

import "bufio"
import "fmt"
import "io"
import "os"
import "path"
import "strings"

func showHelp() {

	fmt.Println("sql-extract <dump.sql> <table>")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("    Extract a specific table and its data from an sql dump file.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    sql-extract dump-1337.sql users_table > users.sql;")
	fmt.Println("")

}

func toFile(chunk string) string {

	var file string

	if strings.HasPrefix(chunk, "/") {

		file = chunk

	} else if strings.HasPrefix(chunk, "./") || strings.HasPrefix(chunk, "../") {

		cwd, err := os.Getwd()

		if err == nil {
			file = path.Join(cwd, chunk)
		}

	} else {

		cwd, err := os.Getwd()

		if err == nil {
			file = path.Join(cwd, chunk)
		}

	}

	return file

}

func main() {

	var file string = ""
	var table string = ""

	if len(os.Args) == 3 {

		file = toFile(os.Args[1])
		table = os.Args[2]

	} else {

		showHelp()
		os.Exit(1)

	}

	stream, err1 := os.Open(file)

	if err1 == nil {

		defer stream.Close()

		reader := bufio.NewReader(stream)

		var matches = false

		for {

			line, err2 := reader.ReadString('\n')

			if err2 != nil {

				if err2 == io.EOF {

					os.Exit(0)

					break

				} else {

					os.Exit(1)
					break

				}

			}

			if strings.Contains(line, "-- Name: "+table+"; Type: TABLE; Schema: public; Owner: -") {
				matches = true
			} else if strings.Contains(line, "-- Data for Name: "+table+";") {
				matches = true
			} else if matches == true {

				if strings.Contains(line, "-- Name:") {
					matches = false
				} else if strings.Contains(line, "-- Data for Name:") {
					matches = false
				}

			}

			if matches == true {
				fmt.Println(line)
			}

		}

	} else {

		os.Exit(1)

	}

}
