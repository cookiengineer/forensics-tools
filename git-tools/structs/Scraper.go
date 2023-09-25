package structs

import "fmt"
import "net/http"
import "io/ioutil"
import "strconv"
import "strings"
import "time"

var CONTENT_TYPES []string = []string{
	"application/gzip",
	"application/json",
	"application/ld+json",
	"application/octet-stream",
	"application/rss+xml",
	"application/x-bzip2",
	"application/x-gzip",
	"application/xml",
	"application/zip",
	"text/html",
	"text/plain",
	"text/xml",
}

type Callback func([]byte)

type ScraperTask struct {
	Url      string
	Callback Callback
}

type Scraper struct {
	Busy      bool
	Limit     int
	Tasks     []ScraperTask
	Headers   map[string]string
	Throttled bool
}

func processRequests(scraper *Scraper) {

	var filtered []ScraperTask
	var limit int = scraper.Limit

	if scraper.Throttled == true {
		limit = 1
	}

	for t := 0; t < len(scraper.Tasks); t++ {

		if len(filtered) < limit {
			filtered = append(filtered, scraper.Tasks[t])
		} else {
			break
		}

	}

	if len(filtered) > 0 {

		for f := 0; f < len(filtered); f++ {

			var task = filtered[f]

			buffer := scraper.Request(task.Url)

			task.Callback(buffer)

		}

		scraper.Tasks = scraper.Tasks[len(filtered):]

		if len(scraper.Tasks) > 0 {

			if scraper.Throttled == true {

				fmt.Println(strconv.Itoa(len(scraper.Tasks)) + " Request Tasks left...")

				time.AfterFunc(10*time.Second, func() {
					processRequests(scraper)
				})

			} else {

				time.AfterFunc(1*time.Second, func() {
					processRequests(scraper)
				})

			}

		} else {

			scraper.Busy = false

		}

	}

}

func NewScraper(limit int) Scraper {

	if limit <= 0 {
		limit = 1
	}

	var scraper Scraper

	scraper.Busy = false
	scraper.Limit = limit
	scraper.Tasks = make([]ScraperTask, 0)
	scraper.Headers = make(map[string]string, 0)
	scraper.Throttled = false

	return scraper

}

func (scraper *Scraper) DeferRequest(url string, callback Callback) {

	scraper.Tasks = append(scraper.Tasks, ScraperTask{
		Url:      url,
		Callback: callback,
	})

	if scraper.Busy == false {

		scraper.Busy = true

		time.AfterFunc(1*time.Second, func() {
			processRequests(scraper)
		})

	}

}

func (scraper *Scraper) Request(url string) []byte {

	var buffer []byte
	var content_type string
	var status_code int

	client := &http.Client{}
	client.CloseIdleConnections()

	request, err1 := http.NewRequest("GET", url, nil)

	if err1 == nil {

		for key, val := range scraper.Headers {
			request.Header.Set(key, val)
		}

		response, err2 := client.Do(request)

		if err2 == nil {

			status_code = response.StatusCode

			if status_code == 200 || status_code == 304 {

				if len(response.Header["Content-Type"]) > 0 {
					content_type = response.Header["Content-Type"][0]
				}

				var valid bool = false

				for c := 0; c < len(CONTENT_TYPES); c++ {

					if strings.Contains(content_type, CONTENT_TYPES[c]) {
						valid = true
						break
					}

				}

				if valid == true {

					data, err3 := ioutil.ReadAll(response.Body)

					if err3 == nil {
						buffer = data
					}

				}

			}

		}

	}

	if len(buffer) > 0 {

		fmt.Println("Request \"" + url + "\"")

	} else {

		fmt.Println("Request \"" + url + "\"")

		if content_type != "" {
			fmt.Println("Unsupported Content-Type \"" + content_type + "\"")
		}

		if status_code != 0 {
			fmt.Println("Unsupported Status Code \"" + strconv.Itoa(status_code) + "\"")
		}

	}

	return buffer

}
