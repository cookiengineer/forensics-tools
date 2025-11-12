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

func showUsage() {

	fmt.Println("Usage: totp-extract <photo-of-qrcode.jpg>")
	fmt.Println("")
	fmt.Println("Usage Notes:")
	fmt.Println("")
	fmt.Println("    This tool will export a photo of a QR code that contains encoded \"otp-migration://\" 2FA seeds.")
	fmt.Println("    Supported 2FA apps are e.g. Google Authenticator, Aegis and others.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    totp-extract DCIM_1337.jpg")
	fmt.Println("")

}

func main() {

	var path string = ""

	if len(os.Args) == 2 {

		if strings.HasSuffix(os.Args[1], ".jpg") {
			path = strings.TrimSpace(os.Args[1])
		} else {

			fmt.Fprintln(os.Stderr, "ERROR: First argument has to be a CRX file")

			showUsage()
			os.Exit(1)

		}

	} else {

		showUsage()
		os.Exit(1)

	}

	if path != "" {

		stat1, err1 := os.Stat(path)

		if err1 == nil && stat1.IsDir() == false {

			file, err2 := os.Open(path)

			if err2 == nil {

				image, _, err3 := image.Decode(file)

				if err3 == nil {

					data := qrcode.Decode(image)

					if strings.HasPrefix(data, "otpauth-migration://offline") {

						fmt.Println("Detected otpauth-migration format")

						results := otpauth.DecodeMigration(data)

						for r, result := range results {

							fmt.Println("> Export authentication issuer \"" + result.Issuer + "\"")

							result_image := qrcode.Encode(*result)
							result_json, _ := json.MarshalIndent(result, "", "\t")

							fmt.Println("> JSON is: " + string(result_json) + "")
							fmt.Println("")

							filename := "totp-qrcode-" + strconv.Itoa(r) + "-" + strings.ToLower(result.Issuer)

							err4 := os.WriteFile(filename+".png", result_image, 0666)

							if err4 == nil {
								fmt.Println("> Exported OTP auth seed to QR Code \"" + filename + ".png\"")
							}

							err5 := os.WriteFile(filename+".json", result_json, 0666)

							if err5 == nil {
								fmt.Println("> Exported OTP auth seed to JSON \"" + filename + ".json\"")
							}

							fmt.Println("")

						}

					} else if strings.HasPrefix(data, "otpauth://") {

						result, err := otpauth.Decode(data)

						if err == nil {

							fmt.Println("> Export authentication issuer \"" + result.Issuer + "\"")

							result_image := qrcode.Encode(*result)
							result_json, _ := json.MarshalIndent(result, "", "\t")

							fmt.Println("> JSON is: " + string(result_json) + "")
							fmt.Println("")

							filename := "totp-qrcode-" + strings.ToLower(result.Issuer)

							err4 := os.WriteFile(filename+".png", result_image, 0666)

							if err4 == nil {
								fmt.Println("> Exported OTP auth seed to QR Code \"" + filename + ".png\"")
							}

							err5 := os.WriteFile(filename+".json", result_json, 0666)

							if err5 == nil {
								fmt.Println("> Exported OTP auth seed to JSON \"" + filename + ".json\"")
							}

							fmt.Println("")

						} else {
							fmt.Fprintln(os.Stderr, "ERROR: Cannot decode OTP auth url")
							fmt.Fprintf(os.Stderr, "ERROR: %v\n", err3)
							os.Exit(1)
						}

					} else {

						fmt.Fprintln(os.Stderr, "ERROR: Unsupported OTP auth format")
						os.Exit(1)

					}

				} else {

					fmt.Fprintln(os.Stderr, "ERROR: Cannot decode JPEG image from \"" + path + "\"")
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err3)
					os.Exit(1)

				}

			} else {

				fmt.Fprintln(os.Stderr, "ERROR: Cannot read from \"" + path + "\"")
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err2)
				os.Exit(1)

			}

		} else {

			fmt.Fprintln(os.Stderr, "ERROR: Cannot read from \"" + path + "\"")
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err1)
			os.Exit(1)

		}

	}

}
