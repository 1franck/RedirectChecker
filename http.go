package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type RequestHeader struct {
	Name  string
	Value string
}

/**
Create a client which don't follow redirect
*/
func createClient() (client *http.Client) {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

/**
Get an url!
*/
func getUrl(client *http.Client, url string, headers []RequestHeader) (response *http.Response, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	for _, v := range headers {
		req.Header.Set(v.Name, v.Value)
	}
	return client.Do(req)
}

/**
Parse a header string (eg: "User-Agent: foobar")
*/
func parseHeader(header string) (RequestHeader, error) {
	parts := strings.Split(header, ":")
	if len(parts[1:]) == 0 {
		return RequestHeader{}, errors.New(fmt.Sprintf("unable to parse header arg \"%v\", missing value part\n", header))
	} else if len(parts) < 2 {
		return RequestHeader{}, errors.New(fmt.Sprintf("unable to parse header arg \"%v\", \"key : val\" format expected\n", header))
	}

	return RequestHeader{
		Name:  strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(fmt.Sprintf("%v", parts[1:])),
	}, nil
}
