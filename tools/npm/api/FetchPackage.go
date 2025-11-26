package api

import "github.com/andybalholm/brotli"
import "compress/flate"
import "compress/gzip"
import "encoding/json"
import "errors"
import "fmt"
import "io"
import "net/http"

func FetchPackage(package_scope string, package_name string, package_version string) (*Schema, error) {

	var schema *Schema

	url := buildURL(package_scope, package_name, package_version)

	if url != "" {

		fmt.Printf("> API URL: \"%s\"\n", url)

		request, err0 := http.NewRequest("GET", url, nil)

		if err0 == nil {

			// request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
			// request.Header.Set("Accept-Encoding", "gzip, deflate, br")

			request.Header.Set("User-Agent", "npm/10.5.0 node/v20.10.0 linux x64")
			request.Header.Set("Accept", "application/vnd.npm.install-v1+json, application/json")
			request.Header.Set("Accept-Encoding", "gzip, deflate, br")

			client := &http.Client{}

			response, err1 := client.Do(request)

			if err1 == nil {

				content_encoding := response.Header.Get("Content-Encoding")

				if response.StatusCode == 200 {

					var reader io.Reader = response.Body

					switch content_encoding {
					case "", "identity":
						reader = response.Body
					case "br":
						reader = brotli.NewReader(response.Body)
					case "gzip":

						gzip_reader, err := gzip.NewReader(response.Body)

						if err != nil {
							return nil, err
						}

						defer gzip_reader.Close()
						reader = gzip_reader

					case "deflate":

						deflate_reader := flate.NewReader(response.Body)
						defer deflate_reader.Close()
						reader = deflate_reader

					}

					bytes, err2 := io.ReadAll(reader)

					if err2 == nil {

						err3 := json.Unmarshal(bytes, &schema)

						if err3 == nil {
							return schema, nil
						} else {
							return nil, err3
						}

					} else {
						return nil, err2
					}

				} else {
					return nil, errors.New(fmt.Sprintf("NPM registry returned HTTP status code %d", response.StatusCode))
				}

			} else {
				return nil, err1
			}

		} else {
			return nil, err0
		}

	}

	return nil, errors.New("Invalid package URL")

}
