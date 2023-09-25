
package main;

import "encoding/binary";
import "fmt";
import "os";
import "strconv";
import "strings";



func showHelp () {

	fmt.Println("");
	fmt.Println("uncrx");
	fmt.Println("");
	fmt.Println("Usage: uncrx ./path/to/extension.crx");
	fmt.Println("");
	fmt.Println("Usage Notes:");
	fmt.Println("");
	fmt.Println("    This tool will create a .zip file in the same folder where the .crx file is");
	fmt.Println("    located, and extract the ZIP data contained within CRX file into it.");
	fmt.Println("");
	fmt.Println("Examples:");
	fmt.Println("");
	fmt.Println("    # creates ~/Downloads/ublock-origin.zip");
	fmt.Println("    uncrx ~/Downloads/ublock-origin.crx");
	fmt.Println("");

}

func main () {

	var buffer []byte;
	var output string;

	if len(os.Args) == 2 {

		if strings.HasSuffix(os.Args[1], ".crx") {

			buf, err := os.ReadFile(strings.TrimSpace(os.Args[1]));

			if err == nil {

				buffer = buf;
				output = os.Args[1][0:strings.LastIndex(os.Args[1], ".crx")] + ".zip";

			}

		} else {

			fmt.Println("First argument has to be a CRX file.");
			os.Exit(1);

		}

	} else {

		showHelp();
		os.Exit(1);

	}


	if len(buffer) > 4 && string(buffer[0]) == "C" && string(buffer[1]) == "r" {

		var header = string(buffer[0]) + string(buffer[1]) + string(buffer[2]) + string(buffer[3]);

		if header == "Cr24" && len(buffer) > 4 {

			if buffer[4] == 0x02 && len(buffer) > 16 {

				fmt.Println("Detected Cr24 / Version 2 file format");

				/*
				 * magic             (4 bytes) Cr24
				 * version           (4 bytes) 2
				 * public key length (4 bytes) length in bytes
				 * signature length  (4 bytes) length in bytes
				 * public key        (*pub key length)
				 * signature         (*signature length)
				 * (... zip data ...)
				 */

				var publickey_length = binary.LittleEndian.Uint32([]byte{
					buffer[8],
					buffer[9],
					buffer[10],
					buffer[11],
				});

				var signature_length = binary.LittleEndian.Uint32([]byte{
					buffer[12],
					buffer[13],
					buffer[14],
					buffer[15],
				});

				var zip_data = buffer[16 + publickey_length + signature_length:];

				fmt.Println("> Public Key length is: " + strconv.Itoa(int(publickey_length)) + " bytes");
				fmt.Println("> Signature length is:  " + strconv.Itoa(int(signature_length)) + " bytes");
				fmt.Println("> ZIP data length is:   " + strconv.Itoa(len(zip_data))         + " bytes");
				fmt.Println("> Write ZIP data to:    \"" + output + "\"");
				fmt.Println("");

				err := os.WriteFile(output, zip_data, 0644);

				if err == nil {

					fmt.Println("Successfully exported ZIP data to \"" + output + "\".");
					os.Exit(0);

				} else {

					fmt.Println("Could not write to \"" + output + "\"!");
					os.Exit(1);

				}

			} else if buffer[4] == 0x03 && len(buffer) > 12 {

				fmt.Println("Detected Cr24 / Version 3 file format");

				/*
				 * Source: https://source.chromium.org/chromium/chromium/src/+/main:components/crx_file/crx3.proto
				 *
				 * magic             (4 bytes) Cr24
				 * version           (4 bytes) 3
				 * crx_header length (4 bytes) little-endian uint32
				 * crx_header        (N bytes)
				 * -> sha256_with_rsa public_key + signature
				 * -> sha256_with_ecdsa public_key + signature
				 * -> verified_contents (4 bytes)
				 * -> signed_header_data (10000)
				 * (... zip_data ...)
				 */

				var crx_header_length = binary.LittleEndian.Uint32([]byte{
					buffer[8],
					buffer[9],
					buffer[10],
					buffer[11],
				});

				var zip_data = buffer[12 + crx_header_length:];

				fmt.Println("> CRX header length is: " + strconv.Itoa(int(crx_header_length)) + " bytes");
				fmt.Println("> ZIP data length is:   " + strconv.Itoa(len(zip_data))          + " bytes");
				fmt.Println("> Write ZIP data to:    \"" + output + "\"");
				fmt.Println("");

				err := os.WriteFile(output, zip_data, 0644);

				if err == nil {

					fmt.Println("Successfully exported ZIP data to \"" + output + "\"!");
					os.Exit(0);

				} else {

					fmt.Println("Could not write to \"" + output + "\"!");
					os.Exit(1);

				}

			}

		} else {

			// TODO: Another Christmas, another future format.

		}

	}

}

