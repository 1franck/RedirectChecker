package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var nbOfJumps = 0
var maxJumps = 20
var redirectTimes []time.Duration

func createClient() (client *http.Client) {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func getUrl(client *http.Client, url string) (response *http.Response, err error) {
	defer timeTrack(time.Now())
	return client.Get(url)
}

func showResponse(response *http.Response) {
	if nbOfJumps > 1 {
		fmt.Printf("\nRedirected to ...\n")
	}
	fmt.Printf("[#%v] %v - %v", nbOfJumps, response.Request.URL.String(), redirectTimes[nbOfJumps-1])
	fmt.Printf("\n > Status: %v\n", response.Status)
	for i, v := range response.Header {
		if i != "Location" {
			fmt.Printf(" > %v: %v\n", i, v)
		}
	}
}

func timeTrack(start time.Time) {
	redirectTimes = append(redirectTimes, time.Since(start))
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Please, specify an url ...")
		return
	}

	url := args[0]
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	client := createClient()

	for {
		resp, err := getUrl(client, url)
		if err != nil {
			fmt.Println("Failed to fetch url:", err)
			break
		}
		resp.Body.Close()

		nbOfJumps++
		showResponse(resp)
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
			fmt.Printf("%v redirects(s) done in %s", nbOfJumps, totalTime)
			break
		}
	}
}
