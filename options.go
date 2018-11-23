package main

import "strings"

type options struct {
	url                 string
	maxJumps            int
	rawRequestHeaders   arrayFlags
	showResponseBody    *bool
	showResponseHeaders *bool
}

/**
Parse header string (ex: "Host: www.example.com" into RequestHeader struct
*/
func (o options) getParsedHeaders() ([]RequestHeader, error) {
	var parsedHeaders []RequestHeader
	for _, v := range o.rawRequestHeaders {
		parsedHeader, err := parseHeader(v)
		if err != nil {
			return parsedHeaders, err
		} else {
			parsedHeaders = append(parsedHeaders, parsedHeader)
		}
	}
	return parsedHeaders, nil
}

/**
Auto prefix the protocol if not present
*/
func (o options) getFormattedUrl() string {
	if !strings.HasPrefix(o.url, "https://") && !strings.HasPrefix(o.url, "http://") {
		return "http://" + o.url
	}
	return o.url
}
