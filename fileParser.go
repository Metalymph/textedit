package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FileParser group all the elements related to the file parser
type FileParser struct {
	inFile, outFile string
	outLineBuilder  *strings.Builder
}

func NewFileParser(inFile, outFile string) *FileParser {
	return &FileParser{inFile: inFile, outFile: outFile}
}

/*
parse open the input and output files and parses input file line by line
producing the output file
*/
func (fp *FileParser) parse() (err error) {
	//reading the input text file as a whole file since there's no specific splitting rule
	inFile, err := os.Open(os.Args[1])
	if err != nil {
		exitWithError(err.Error())
	}
	defer inFile.Close()

	//opening the destination file
	outFile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		exitWithError(err.Error())
	}
	defer outFile.Close()

	//building work buffer scanner
	fileScanner := bufio.NewScanner(inFile)
	fileScanner.Split(bufio.ScanLines)

	fp.outLineBuilder = &strings.Builder{}

	//read file line by line and change
	for fileScanner.Scan() {
		if err := fileScanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not scan the file due to this %s error \n", err.Error())
			os.Exit(3)
		}

		err := fp.parseLine(fileScanner.Text())
		if err != nil {
			return err
		}

		newLine := fp.outLineBuilder.String()
		// changes the last space char with the newline to save the line format
		adjustedNewLine := newLine[:len(newLine)-1] + "\n"
		_, err = outFile.WriteString(adjustedNewLine)
		if err != nil {
			return err
		}
		fp.outLineBuilder.Reset()
	}
	return nil
}

// parseLine parses a single line and edits it following the given rules
func (fp *FileParser) parseLine(line string) (err error) {
	fp.outLineBuilder.Grow(len(line) + 1)
	tokens := strings.Split(line, " ")

	// this is useful to know the last punct. found (true -> open, false - closed)
	punctMarkStatus := false

	for i, currentToken := range tokens {
		switch currentToken[0] {
		case '.', ',', '!', '?', ':', ';', '\'':
			if err = fp.outLineBuilder.WriteByte(currentToken[0]); err != nil {
				return err
			}
			if len(currentToken) > 1 {
				//example ",a angel" becomes "a angel" to check below
				//example "'a angel" becomes "a angel" to check below
				currentToken = currentToken[1:]
			} else {
				continue
			}
		default:
		}

		// write punct. until a different char is found
		for i := range currentToken {
			switch currentToken[i] {
			case '.', ',', '!', '?', ':', ';':
				if err = fp.outLineBuilder.WriteByte(currentToken[i]); err != nil {
					return err
				}
			default:
				currentToken = currentToken[i+1:]
			}
		}

		// open and close ' cases
		if currentToken[0] == '\'' {
			if err = fp.outLineBuilder.WriteByte(currentToken[0]); err != nil {
				return err
			}
			// a whitespace must be written only if this is a closing '
			if punctMarkStatus {
				if err = fp.outLineBuilder.WriteByte(' '); err != nil {
					return err
				}
			}
			punctMarkStatus = !punctMarkStatus
			if len(currentToken) > 1 {
				//example "'a angel" becomes "a angel" to check below
				currentToken = currentToken[1:]
			} else {
				continue
			}
		}

		switch currentToken {
		// previous token verification cases
		case "(cap)":
			capitalizedStr := strings.ToUpper(tokens[i-1][:1]) + tokens[i-1][1:]
			if err := fp.appendToLine(capitalizedStr); err != nil {
				return err
			}
		case "(low)":
			if err := fp.appendToLine(strings.ToLower(tokens[i-1])); err != nil {
				return err
			}
		case "(up)":
			if err := fp.appendToLine(strings.ToUpper(tokens[i-1])); err != nil {
				return err
			}
		case "(bin)":
			decimalNum, err := strconv.ParseInt(tokens[i-1], 2, 64)
			if err != nil {
				return err
			}
			if err = fp.appendToLine(strconv.FormatInt(decimalNum, 10)); err != nil {
				return err
			}
		case "(hex)":
			decimalNum, err := strconv.ParseInt(tokens[i-1], 16, 64)
			if err != nil {
				return err
			}
			if err = fp.appendToLine(strconv.FormatInt(decimalNum, 10)); err != nil {
				return err
			}

		//next token verification cases
		case "a", "A":
			//we can't know if 'a' is the last token of the input file string
			if i+1 < len(tokens) {
				//if the next token first char is a vowel than writes 'an' or 'An'
				switch tokens[i+1][0] {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					if err := fp.appendToLine(currentToken + "n"); err != nil {
						return err
					}
				}
			}
		case "(cap,":
			numOcc, err := strconv.ParseUint(tokens[i+1], 10, 64)
			if err != nil {
				return err
			}
			capitalizeStrings(tokens[i-int(numOcc) : i])
		case "(low,":
			numOcc, err := strconv.ParseUint(tokens[i+1], 10, 64)
			if err != nil {
				return err
			}
			lowerStrings(tokens[i-int(numOcc) : i])
		case "(up,":
			numOcc, err := strconv.ParseUint(tokens[i+1], 10, 64)
			if err != nil {
				return err
			}
			upperStrings(tokens[i-int(numOcc) : i])
		default:
		}
	}
	return nil
}

// appendToLine writes a new piece of tthe new line on the builder
func (fp *FileParser) appendToLine(chunk string) (err error) {
	_, err = fp.outLineBuilder.WriteString(chunk)
	if err != nil {
		return err
	}
	err = fp.outLineBuilder.WriteByte(' ')
	if err != nil {
		return err
	}
	return nil
}