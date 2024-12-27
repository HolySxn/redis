package main

import (
	"strings"
)

func isValidPort(port int) bool {
	if port < 1024 || port > 49151 {
		return false
	}

	return true
}

func PrepareInput(input []byte) []string {
	s := strings.TrimSuffix(string(input), "\n")
	return strings.Split(s, " ")
}
