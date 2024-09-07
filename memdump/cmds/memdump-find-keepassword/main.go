package main

import "fmt"
import "log"
import "os"

func main() {

	if len(os.Args) == 2 {

		buffer, err := os.ReadFile(os.Args[1])

		if err != nil {
			log.Fatal(err)
		}

		candidates := make([]int, 0)

		for b := 0; b < len(buffer); b++ {

			var offset = b

			if buffer[offset] == 0xCF && buffer[offset + 1] == 0x25 {
				candidates = append(candidates, offset)
			}

		}

		for c := 0; c < len(candidates); c++ {

			offset := candidates[c] + 2

			fmt.Println("Candidate offset is: ", offset)

			// assume less than 32 bytes password length
			for length := 0; length < 32; length++ {

				// 0x20 - 0xFF are valid characters
				if buffer[offset + length] >= 0x20 && buffer[offset + length] <= 0xFF {
					fmt.Println("Password might be: ", length, string(buffer[offset:offset + length]))
				} else {
					break
				}

			}

		}

	}

}
