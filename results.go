package main

import (
	"fmt"
	"net/http"
)

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
