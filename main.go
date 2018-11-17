package main

import (
	"fmt"
	"net/http"
	"os"
)

var nbOfJumps = 0

func createClient() (client *http.Client) {
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

func showResponse(response *http.Response) {
	if nbOfJumps > 1 {
		fmt.Printf("\nRedirected to ...\n")
	}
	fmt.Printf("[#%v] %v", nbOfJumps, response.Request.URL.String())
	fmt.Printf("\n > Status: %v\n", response.Status)
	for i, v := range response.Header {
		if i != "Location" {
			fmt.Printf(" > %v: %v\n", i, v)
		}
	}
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Please, specify an url ...")
		return
	}

	url := args[0]
	client := createClient()

	for {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Failed to fetch url:", err)
			break
		}
		resp.Body.Close()

		nbOfJumps++
		showResponse(resp)
		if resp.Header.Get("location") != "" {
			url = resp.Header.Get("location")
		} else {
			break
		}
	}

	nbOfJumps = 0
}
