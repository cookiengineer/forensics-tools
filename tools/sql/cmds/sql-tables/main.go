package main

import "bufio"
import "fmt"
import "io"
import "os"
import "strings"

func showUsage() {

	fmt.Println("Usage: sql-tables <dump.sql>")
	fmt.Println("")
	fmt.Println("    This tools extracts the table index from a large SQL dump file.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    sql-tables theleak-1337.sql > list-of-tables.txt;")
	fmt.Println("")

}

func main() {

	var sql_file string

	if len(os.Args) == 2 {

		if strings.HasSuffix(os.Args[1], ".sql") {
			sql_file = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be an SQL file")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if sql_file != "" {

		stream, err0 := os.Open(sql_file)

		if err0 == nil {

			reader := bufio.NewReader(stream)
			defer stream.Close()

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

				if strings.Contains(line, "CREATE TABLE public.") {

					chunk := strings.TrimSpace(line[20:])

					if strings.Contains(chunk, "(") && strings.Contains(chunk, ")") {
						chunk = strings.TrimSpace(chunk[0:strings.Index(chunk, "(")])
					} else if strings.Contains(chunk, "(") {
						chunk = strings.TrimSpace(chunk[0:strings.Index(chunk, "(")])
					}

					fmt.Println(chunk)

				}

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read from \"" + sql_file + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err0)
			os.Exit(1)

		}

	}

}
