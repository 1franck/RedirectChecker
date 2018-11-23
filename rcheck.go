package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var nbOfJumps = 0
var redirectTimes []time.Duration

func timeTrack(start time.Time) {
	redirectTimes = append(redirectTimes, time.Since(start))
}

func httpGet(client *http.Client, url string, headers []RequestHeader) (response *http.Response, err error) {
	defer timeTrack(time.Now())
	return getUrl(client, url, headers)
}

func main() {

	options := options{
		url:                 "",
		maxJumps:            20,
		rawRequestHeaders:   nil,
		showResponseHeaders: flag.Bool("i", false, "Show response headers"),
		showResponseBody:    flag.Bool("b", false, "Show response body"),
	}

	flag.Var(&options.rawRequestHeaders, "H", "Header for the request. Flag can be repeat for multiple header values")
	flag.Parse()
	options.url = flag.Arg(0)

	if len(options.url) < 1 {
		fmt.Println("Please, specify an url ...")
		return
	}

	parsedHeaders, err := options.getParsedHeaders()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := createClient()
	url := options.getFormattedUrl()

	for {
		resp, err := httpGet(client, url, parsedHeaders)
		if err != nil {
			fmt.Println("Failed to fetch url:", err)
			break
		}

		nbOfJumps++
		showResponse(resp, *options.showResponseHeaders)

		if *options.showResponseBody {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}
		resp.Body.Close()

		if resp.Header.Get("location") != "" {
			if nbOfJumps > options.maxJumps {
				fmt.Printf("Maximum of %v redirects reached!", options.maxJumps)
				break
			}
			url = resp.Header.Get("location")
		} else {
			var totalTime time.Duration
			for _, v := range redirectTimes {
				totalTime += v
			}
			fmt.Printf("%v redirects(s) done in %s\n", nbOfJumps, totalTime)
			break
		}
	}
}
