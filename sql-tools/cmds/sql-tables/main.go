
package main;

import "bufio";
import "fmt";
import "io";
import "os";
import "path";
import "strings";



func showHelp () {

	fmt.Println("sql-tables <dump.sql>");
	fmt.Println("");
	fmt.Println("Usage:");
	fmt.Println("");
	fmt.Println("    Extract the table index from an sql dump file.");
	fmt.Println("");
	fmt.Println("Examples:");
	fmt.Println("");
	fmt.Println("    sql-tables dump-1337.sql > tables.txt;");
	fmt.Println("");

}

func toFile (chunk string) string {

	var file string;

	if strings.HasPrefix(chunk, "/") {

		file = chunk;

	} else if strings.HasPrefix(chunk, "./") || strings.HasPrefix(chunk, "../") {

		cwd, err := os.Getwd();

		if err == nil {
			file = path.Join(cwd, chunk);
		}

	} else {

		cwd, err := os.Getwd();

		if err == nil {
			file = path.Join(cwd, chunk);
		}

	}

	return file;

}



func main () {

	var file string = "";

	if len(os.Args) == 2 {

		file = toFile(os.Args[1]);

	} else {

		showHelp();
		os.Exit(1);

	}


	stream, err1 := os.Open(file);

	if err1 == nil {

		defer stream.Close();

		reader := bufio.NewReader(stream);


		for {

			line, err2 := reader.ReadString('\n');

			if err2 != nil {

				if err2 == io.EOF {

					os.Exit(0);

					break;

				} else {

					fmt.Println(err2);

					os.Exit(1);
					break;

				}

			}


			if strings.Contains(line, "CREATE TABLE public.") {

				var chunk = strings.TrimSpace(strings.Split(line, "CREATE TABLE public.")[1]);

				if strings.Contains(chunk, "(") {
					chunk = strings.TrimSpace(strings.Split(chunk, "(")[0]);
				}

				fmt.Println(chunk);

			}

		}

	} else {
		os.Exit(1);
	}

}

