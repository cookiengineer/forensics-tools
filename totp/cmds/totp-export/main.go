package main

import "totp/otpauth"
import "totp/qrcode"
import "encoding/json"
import "fmt"
import "image"
import _ "image/jpeg"
import "os"
import "strconv"
import "strings"

func main() {

	var path string = ""

	if len(os.Args) == 2 {
		path = strings.TrimSpace(os.Args[1])
	}

	if path != "" {

		stat1, err1 := os.Stat(path)

		if err1 == nil && stat1.IsDir() == false {

			file, err2 := os.Open(path)
			image, _, err3 := image.Decode(file)

			if err2 == nil && err3 == nil {

				data := qrcode.Decode(image)

				if strings.HasPrefix(data, "otpauth://") {

					// Do Nothing
					fmt.Println("Currently unsupported!")
					fmt.Println(data)

				} else if strings.HasPrefix(data, "otpauth-migration://offline") {

					results := otpauth.DecodeMigration(data)

					for r := 0; r < len(results); r++ {

						result := results[r]
						result_image := qrcode.Encode(result)
						result_json, _ := json.MarshalIndent(result, "", "\t")

						prefix := "qrcode-" + strconv.Itoa(r) + "-" + strings.ToLower(result.Issuer)

						os.WriteFile("/tmp/"+prefix+".png", result_image, 0666)
						os.WriteFile("/tmp/"+prefix+".json", result_json, 0666)

						fmt.Println("Exported \"/tmp/" + prefix + ".{png,json}\" from \"" + result.Name + "\" (" + result.Issuer + ")")

					}

				}

			} else {
				os.Exit(1)
			}

		} else {
			os.Exit(1)
		}

	} else {

		fmt.Println("TOTP Extractor")
		fmt.Println("")
		fmt.Println("Usage: totp-extract ./path/to/camera-photo-of-qrcode.jpg")
		fmt.Println("")

		os.Exit(2)

	}

}
