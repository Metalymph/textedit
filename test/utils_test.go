package test

import (
	"testing"
	fp "github.com/Metalymph/textedit/fileparser"
)


func TestCapitalizeStrings(t *testing.T) {
	strs := []string{"ok", "dream", "please"}
	capStrs := [3]string{"Ok", "Dream", "Please"}

	fp.CapitalizeStrings(strs)
	for i := range strs {
		if  strs[i] != capStrs[i] {
			t.Fatalf("%s -> must be -> %s", strs[i], capStrs[i])
		}
	}
}

func TestLowerStrings(t *testing.T) {
	strs := []string{"OK", "DREAM", "PLEASE"}
	lowStrs := [3]string{"ok", "dream", "please"}

	fp.LowerStrings(strs)
	for i := range strs {
		if  strs[i] != lowStrs[i] {
			t.Fatalf("%s -> must be -> %s", strs[i], lowStrs[i])
		}
	}
}

func TestUpperStrings(t *testing.T) {
	strs := []string{"ok", "dream", "please"}
	upStrs := [3]string{"OK", "DREAM", "PLEASE"}

	fp.UpperStrings(strs)
	for i := range strs {
		if  strs[i] != upStrs[i] {
			t.Fatalf("%s -> must be -> %s", strs[i], upStrs[i])
		}
	}
}