package main

import (
	"fmt"
	"os"
	"strings"
)

// UTILS
func capitalizeStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToUpper(strs[i-1][:1]) + strs[i-1][1:]
	}
}

func lowerStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToLower(strs[i])
	}
}

func upperStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToUpper(strs[i])
	}
}

func exitWithError(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}