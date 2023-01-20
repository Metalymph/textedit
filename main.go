package main

import (
	"fmt"
	"os"
	fp "github.com/Metalymph/textedit/fileparser"
)

func main() {
	//number of arguments control
	argsLen := len(os.Args)
	if argsLen != 3 {
		fp.ExitWithError(fmt.Sprintf("Wrong number of arguments. 2 wanted, %d given", argsLen))
	}

	if err := fp.NewFileParser(os.Args[1], os.Args[1]).Parse(); err != nil {
		fp.ExitWithError(err.Error())
	}
}
