package jacktokenizer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var (
	lineNumber       int
	pos              int
	bfr              *bufio.Scanner
	currenLine       string
	currentToken     string
	currentTokenList []string
	currentCFile     string

	file *os.File

	inBlockComment, inString, inCurrly, inBracket bool // These are the parser state
)

func SetCurrentCFile(file string) {
	currentCFile = file
}

func GetCurrentTokensList() []string {
	return currentTokenList
}

func GetCurrentLineNumber() int {
	return lineNumber
}

func GetTokenType(token string) string {
	if _, ok := Keywords[token]; ok {
		return "keyword"
	} else if _, ok := Symbols[token]; ok {
		return "symbol"
	} else if _, err := strconv.Atoi(token); err == nil {
		return "integerConstant"
	} else if strings.HasPrefix(token, "\"") {
		return "stringConstant"
	} else {
		return "identifier"
	}
}

func OpenFile(pathName string) {
	lineNumber = 1

	inputInfo, err := os.Stat(pathName)
	handleError(err)

	if !inputInfo.IsDir() {
		file, err = os.Open(pathName)
		handleError(err)

		bfr = bufio.NewScanner(file)
		bfr.Split(bufio.ScanLines)
	}
}

func handleTokens(currentLine string) {
	tempToken := ""
	currentTokenList = make([]string, 0)
	inLineComment := false

	for i, letter := range currentLine {
		tempToken += string(letter)
		pos = i

		if letter == '"' {
			if inString {
				inString = false
				currentToken = tempToken
				tempToken = ""
			} else {
				inString = true
			}
		}
		if letter == ' ' && !inString {
			currentToken = tempToken
			tempToken = ""
		} else {
			if letter == '/' && i < len(currentLine)-1 && currentLine[i+1] == '/' {
				inLineComment = true
			}
			if inLineComment {
				tempToken = ""
				continue
			}
			if letter == '*' && i < len(currentLine)-1 && currentLine[i+1] == '/' {
				tempToken = ""
				if !inBlockComment {
					fmt.Println(currentCFile, "Missing /* on position ", lineNumber, ", ", pos)
				} else {
					inBlockComment = false
					continue
				}
			} else if letter == '/' && i < len(currentLine)-1 && currentLine[i+1] == '*' {
				tempToken = ""
				if inBlockComment {
					fmt.Println(currentCFile, "Missing */ on position ", lineNumber, ", ", pos)
					tempToken = ""
				} else {
					inBlockComment = true
				}
			} else if letter == '/' && len(currentLine) > 1 && currentLine[i-1] == '*' {
				tempToken = ""
				inBlockComment = false
			}

			if inBlockComment {
				tempToken = ""
				continue
			}
			if !inString {
				if _, ok := Keywords[tempToken]; ok && !unicode.IsLetter(rune(currentLine[i+1])) {
					currentToken = tempToken
					tempToken = ""
				} else if _, ok := Symbols[tempToken]; ok {
					currentToken = tempToken
					tempToken = ""
				} else if letter == '{' {
					if inCurrly {
						fmt.Println(currentCFile, "Missing } on line", lineNumber)
					} else {
						inCurrly = true
					}
				} else if letter == '[' {
					if inBracket {
						fmt.Println(currentCFile, "Missing ] on line", lineNumber)
					} else {
						inBracket = true
					}
				} else if unicode.IsDigit(letter) {
					if i < len(currentLine)-1 && !unicode.IsDigit(rune(currentLine[i+1])) {
						currentToken = tempToken
						tempToken = ""
					}
				} else if unicode.IsLetter(letter) || letter == '_' {
					if i < len(currentLine)-1 && (!unicode.IsLetter(rune(currentLine[i+1])) && !unicode.IsDigit(rune(currentLine[i+1])) && currentLine[i+1] != '_') {
						currentToken = tempToken
						tempToken = ""
					}
				} else if letter == '}' {
					currentToken = string(letter)
					tempToken = ""

					if !inCurrly {
						fmt.Println(currentCFile, "Missing { on line", lineNumber)
					} else {
						inCurrly = false
					}
				} else if letter == ']' {
					currentToken = string(letter)
					tempToken = ""

					if !inBracket {
						fmt.Println(currentCFile, "Missing [ on line", lineNumber)
					} else {
						inBracket = false
					}
				}
			} else if string(letter) == "\n" && inString {
				fmt.Println(currentCFile, "Missing \" in line", lineNumber)
			}
		}
		if currentToken != "" && currentToken != " " {
			currentTokenList = append(currentTokenList, currentToken)
		}
		currentToken = ""
	}
}

func Advance() bool {
	if bfr.Scan() {
		lineNumber++
		currenLine = strings.TrimSpace(bfr.Text())

		if strings.HasPrefix(currenLine, "//") {
			bfr.Scan()
			lineNumber++
			currenLine = bfr.Text()
		}

		if !isValidParentheses(currenLine) {
			fmt.Println(currentCFile, "Invalid parentheses on line ", lineNumber)
		}

		handleTokens(currenLine)
	} else {
		return false
	}
	return true
}

func isValidParentheses(parens string) bool {
	n := 0
	for i := 0; i < len(parens); i++ {
		if parens[i] == '(' {
			n++
		}
		if parens[i] == ')' {
			n--
		}
		if n < 0 {
			return false
		}
	}

	return n == 0
}

func CloseFile() {
	defer file.Close()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
