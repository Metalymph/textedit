package main

import (
	"log"
	"os"
	fp "github.com/Metalymph/textedit/fileparser"
)

func main() {
	//number of arguments control
	argsLen := len(os.Args)
	if argsLen != 3 {
		log.Fatalf("Wrong number of arguments. 2 wanted, %d given", argsLen-1)
	}

	if err := fp.NewFileParser(os.Args[1], os.Args[1]).Parse(); err != nil {
		log.Fatal(err.Error())
	}
}
