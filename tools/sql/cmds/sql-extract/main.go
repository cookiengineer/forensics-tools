package main

import "bufio"
import "fmt"
import "io"
import "os"
import "strings"

func showUsage() {

	fmt.Println("Usage: sql-extract <dump.sql> <table_name>")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("    This tool extracts a specific table and its data from a large SQL dump file.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    sql-extract theleak-1337.sql users_table > users.sql;")
	fmt.Println("")

}

func main() {

	var sql_file string
	var sql_table string

	if len(os.Args) == 3 {

		if strings.HasSuffix(os.Args[1], ".sql") {
			sql_file = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be an SQL file")

			showUsage()
			os.Exit(1)

		}

		if strings.TrimSpace(os.Args[2]) != "" {
			sql_table = strings.TrimSpace(os.Args[2])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Second argument has to be an SQL table name")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}


	if sql_file != "" && sql_table != "" {

		stream, err0 := os.Open(sql_file)

		if err0 == nil {

			reader := bufio.NewReader(stream)
			defer stream.Close()

			matches := false

			for {

				line, err1 := reader.ReadString('\n')

				if err1 != nil {

					if err1 == io.EOF {

						os.Exit(0)
						break

					} else {

						os.Exit(1)
						break

					}

				}

				if strings.Contains(line, "-- Name: "+sql_table+"; Type: TABLE; Schema: public; Owner: -") {
					matches = true
				} else if strings.Contains(line, "-- Data for Name: "+sql_table+";") {
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

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read from \"" + sql_file + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
