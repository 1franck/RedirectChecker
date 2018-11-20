package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var nbOfJumps = 0
var maxJumps = 20
var redirectTimes []time.Duration

func showResponse(response *http.Response, showHeaders bool) {
	if nbOfJumps > 1 {
		fmt.Printf("\nRedirected to ...\n")
	}
	fmt.Printf("[#%v] %v - %v", nbOfJumps, response.Request.URL.String(), redirectTimes[nbOfJumps-1])
	fmt.Printf("\n > Status: %v\n > Protocol: %v\n", response.Status, response.Proto)
	if showHeaders {
		showResponseHeaders(response)
	}
}

func showResponseHeaders(response *http.Response) {
	for i, v := range response.Header {
		if i != "Location" {
			fmt.Printf(" > %v: %v\n", i, v)
		}
	}
}

func timeTrack(start time.Time) {
	redirectTimes = append(redirectTimes, time.Since(start))
}

func httpGet(client *http.Client, url string, headers []RequestHeader) (response *http.Response, err error) {
	defer timeTrack(time.Now())
	return getUrl(client, url, headers)
}

func parseHeadersFlag(headersFlag arrayFlags) ([]RequestHeader, error) {
	var parsedHeaders []RequestHeader
	for _, v := range headersFlag {
		parsedHeader, err := parseHeader(v)
		if err != nil {
			return parsedHeaders, err
		} else {
			parsedHeaders = append(parsedHeaders, parsedHeader)
		}
	}
	return parsedHeaders, nil
}

func main() {

	inspectHeadersFlag := flag.Bool("i", false, "Show response headers")
	showBodyFlag := flag.Bool("b", false, "Show response body")

	var requestHeaderFlag arrayFlags
	flag.Var(&requestHeaderFlag, "H", "Header for the request. Flag can be repeat for multiple header values")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please, specify an url ...")
		return
	}
	url := args[0]

	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	parsedHeaders, err := parseHeadersFlag(requestHeaderFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := createClient()

	for {
		resp, err := httpGet(client, url, parsedHeaders)
		if err != nil {
			fmt.Println("Failed to fetch url:", err)
			break
		}

		nbOfJumps++
		showResponse(resp, *inspectHeadersFlag)

		if *showBodyFlag {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}
		resp.Body.Close()

		if resp.Header.Get("location") != "" {
			if nbOfJumps > maxJumps {
				fmt.Printf("Maximum of %v redirects reached!", maxJumps)
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
