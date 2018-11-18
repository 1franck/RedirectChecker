package main

import "strings"

/**
Allow multiple values for a cli flag
*/
type arrayFlags []string

func (i *arrayFlags) String() string {
	// this is just an example to satisfy the interface
	var stringRep string
	for _, v := range *i {
		stringRep += v
	}
	return stringRep
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

//func timeTrack(start time.Time) {
//	redirectTimes = append(redirectTimes, time.Since(start))
//}
