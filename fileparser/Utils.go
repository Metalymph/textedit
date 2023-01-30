package fileparser

import (
	"strings"
)

// UTILS
func CapitalizeStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToUpper(strs[i][:1]) + strs[i][1:]
	}
}

func LowerStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToLower(strs[i])
	}
}

func UpperStrings(strs []string) {
	for i := range strs {
		strs[i] = strings.ToUpper(strs[i])
	}
}