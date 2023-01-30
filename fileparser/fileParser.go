package fileparser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// FileParser group all the elements related to the file parser
type FileParser struct {
	inFile, outFile string
	fixedWords      []string
	outLineBuilder  *strings.Builder
}

func NewFileParser(inFile, outFile string) *FileParser {
	return &FileParser{
		inFile:     inFile,
		outFile:    outFile,
		fixedWords: make([]string, 0),
	}
}

/*
parse open the input and output files and parses input file line by line
producing the output file
*/
func (fp *FileParser) Parse() (err error) {
	//reading the input text file
	inFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer inFile.Close()

	//opening the destination file
	outFile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outFile.Close()

	//building work buffer scanner
	fileScanner := bufio.NewScanner(inFile)
	fileScanner.Split(bufio.ScanLines)

	fp.outLineBuilder = &strings.Builder{}

	//read input file line by line
	for fileScanner.Scan() {
		if err := fileScanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not scan the file due to this %s error \n", err.Error())
			os.Exit(3)
		}

		line := fileScanner.Text()
		fp.outLineBuilder.Grow(len(line) + 1) //adjust row memory necessary

		//modify the read line with the given rules
		err := fp.parseLine(line)
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

func (fp *FileParser) parsePunc(currentToken string) (string, error) {
	// write punct. until a different char is found
	var j int
	tokenLen := len(currentToken)
out:
	for j < tokenLen {
		switch currentToken[j] {
		case '.', ',', '!', '?', ':', ';':
		default:
			break out
		}
		j++
	}

	// write a space only if at least one punct is found and proceed with the word's rest
	// for instance: ",time" becomes "," and "time".
	if j > 0 {
		// The "," must be appended to the previous word.
		fp.fixedWords[len(fp.fixedWords)-1] += currentToken[:j]
	}

	if j == tokenLen {
		// only punct. composed token
		return "", nil
	} else {
		// returns the rest of the token: "time" for our example
		return currentToken[j:], nil
	}
}

// parseLine parses a single line and edits it following the given rules
func (fp *FileParser) parseLine(line string) (err error) {
	fmt.Println("**LINE ->", line)
	tokens := strings.Split(line, " ")
	// this is useful to know the last punct. found (true -> open, false - closed) for 'text' cases
	punctMarkStatus := false

	for i, currentToken := range tokens {
		fmt.Println("fixed words ->", fp.fixedWords)
		fmt.Println("\tnext token ->", currentToken)

		// skip n) cases because they're related to (cap, (low, and (up,  cases
		if len(currentToken) == 2 && currentToken[1] == ')' {
			fmt.Println("\t\t-> Token Skipped")
			continue
		}

		currentToken, err = fp.parsePunc(currentToken)
		if err != nil {
			return err
		}
		// the token contained only punct. symbols so we can proceed to the next token
		if currentToken == "" {
			continue
		}
		fmt.Printf("\tafter punct & ' check -> %s\n", currentToken)

		// open ' cases
		if currentToken[0] == '\'' {
			if !punctMarkStatus {
				//it's the first '
				fp.fixedWords = append(fp.fixedWords, currentToken)
				punctMarkStatus = !punctMarkStatus
				continue
			} else {
				//it's the second ' so must be appended to the previous word
				fp.fixedWords[len(fp.fixedWords)-1] += currentToken[:1]
			}
			punctMarkStatus = !punctMarkStatus

			if len(currentToken) > 1 {
				//example "'a angel" becomes "a angel" to check below
				currentToken = currentToken[1:]
			} else {
				// only ' found
				continue
			}
		}

		curTokenLen := len(currentToken)
		//close ' case, when is an open ' (for instance: a angel')
		if currentToken[curTokenLen-1] == '\'' && !punctMarkStatus {
			// becomes angel '
			tokenFixed := fmt.Sprintf("%s %b", currentToken[:curTokenLen-1], currentToken[curTokenLen-1])
			fp.fixedWords = append(fp.fixedWords, tokenFixed)
			punctMarkStatus = !punctMarkStatus
		}

		fixedWordLen := len(fp.fixedWords)
		// at this point the originale token could have been modified with punct and ' checks
		switch currentToken {
		// previous token verification cases
		case "(cap)":
			fp.fixedWords[i-1] = strings.ToUpper(fp.fixedWords[fixedWordLen-1][:1]) + fp.fixedWords[fixedWordLen-1][1:]
		case "(low)":
			fp.fixedWords[fixedWordLen-1] = strings.ToLower(fp.fixedWords[fixedWordLen-1])
		case "(up)":
			fp.fixedWords[fixedWordLen-1] = strings.ToUpper(fp.fixedWords[fixedWordLen-1])
		case "(bin)":
			decimalNum, err := strconv.ParseInt(fp.fixedWords[fixedWordLen-1], 2, 64)
			if err != nil {
				return err
			}
			fp.fixedWords[fixedWordLen-1] = strconv.FormatInt(decimalNum, 10)
		case "(hex)":
			decimalNum, err := strconv.ParseInt(fp.fixedWords[fixedWordLen-1], 16, 64)
			if err != nil {
				return err
			}
			fp.fixedWords[fixedWordLen-1] = strconv.FormatInt(decimalNum, 10)
		//next token verification cases
		case "(cap,":
			numOcc, err := strconv.ParseUint(tokens[i+1][:1], 10, 64)
			if err != nil {
				return err
			}
			CapitalizeStrings(fp.fixedWords[fixedWordLen-int(numOcc):])
		case "(low,":
			numOcc, err := strconv.ParseUint(tokens[i+1][:1], 10, 64)
			if err != nil {
				return err
			}
			LowerStrings(fp.fixedWords[fixedWordLen-int(numOcc):])
		case "(up,":
			numOcc, err := strconv.ParseUint(tokens[i+1][:1], 10, 64)
			if err != nil {
				return err
			}
			UpperStrings(fp.fixedWords[fixedWordLen-int(numOcc):])
		case "a", "A":
			// if 'a' isn't the last token of the input file string
			if i+1 < len(tokens) {
				//if the next token first char is a vowel than writes 'an' or 'An'
				switch tokens[i+1][0] {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					fp.fixedWords = append(fp.fixedWords, currentToken+"n")
				}
			}
		default:
			/*
				here all words and numbers are skipped (appended to fp.fixedWords) and written
				only after we meet a type of token over here.
			*/
			fp.fixedWords = append(fp.fixedWords, currentToken)
		}
	}

	// writes modified line into the string builder
	fp.buildModifiedLine()

	return nil
}

// appendToLine writes a new piece of the new line on the builder
func (fp *FileParser) buildModifiedLine() (err error) {
	for i, newLineChunk := range fp.fixedWords {
		_, err = fp.outLineBuilder.WriteString(newLineChunk)
		if err != nil {
			return err
		}
		// for the last fixed word a following space is useless
		if i == len(fp.fixedWords)-1 {
			continue
		}
		err = fp.outLineBuilder.WriteByte(' ')
		if err != nil {
			return err
		}
	}
	return nil
}
