package main

import (
	"fmt"
	"os"
)

func main() {
	//number of arguments control
	argsLen := len(os.Args)
	if argsLen != 3 {
		exitWithError(fmt.Sprintf("Wrong number of arguments. 2 wanted, %d given", argsLen))
	}

	if err := NewFileParser(os.Args[1], os.Args[1]).parse(); err != nil {
		exitWithError(err.Error())
	}
}
