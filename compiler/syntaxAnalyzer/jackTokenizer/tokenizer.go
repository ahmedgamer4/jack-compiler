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

	file *os.File

	inBlockComment, inString, inCurrly, inBracket bool // These are the parser state

	keywords = map[string]int{
		"class":       0,
		"function":    0,
		"constructor": 0,
		"int":         0,
		"var":         0,
		"boolean":     0,
		"char":        0,
		"void":        0,
		"static":      0,
		"field":       0,
		"let":         0,
		"if":          0,
		"else":        0,
		"while":       0,
		"return":      0,
		"true":        0,
		"false":       0,
		"null":        0,
		"this":        0,
		"do":          0,
		"method":      0,
	}

	symbols = map[string]int{
		"{": 0,
		"}": 0,
		"(": 0,
		")": 0,
		"[": 0,
		"]": 0,
		".": 0,
		",": 0,
		";": 0,
		"+": 0,
		"-": 0,
		"*": 0,
		"/": 0,
		"&": 0,
		"|": 0,
		"<": 0,
		">": 0,
		"=": 0,
		"~": 0,
	}
)

func GetCurrentTokensList() []string {
	return currentTokenList
}

func GetCurrentLineNumber() int {
	return lineNumber
}

func GetTokenType(token string) string {
	if _, ok := keywords[token]; ok {
		return "keyword"
	} else if _, ok := symbols[token]; ok {
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
					fmt.Println("Missing /* on position ", lineNumber, ", ", pos)
				} else {
					inBlockComment = false
					continue
				}
			} else if letter == '/' && i < len(currentLine)-1 && currentLine[i+1] == '*' {
				tempToken = ""
				if inBlockComment {
					fmt.Println("Missing */ on position ", lineNumber, ", ", pos)
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
				if _, ok := keywords[tempToken]; ok {
					currentToken = tempToken
					tempToken = ""
				} else if _, ok := symbols[tempToken]; ok {
					currentToken = tempToken
					tempToken = ""
				} else if letter == '{' {
					if inCurrly {
						fmt.Println("Missing } on line", lineNumber)
					} else {
						inCurrly = true
					}
				} else if letter == '[' {
					if inBracket {
						fmt.Println("Missing ] on line", lineNumber)
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
						fmt.Println("Missing { on line", lineNumber)
					} else {
						inCurrly = false
					}
				} else if letter == ']' {
					currentToken = string(letter)
					tempToken = ""

					if !inBracket {
						fmt.Println("Missing [ on line", lineNumber)
					} else {
						inBracket = false
					}
				}
			} else if string(letter) == "\n" && inString {
				fmt.Println("Missing \" in line", lineNumber)
			}
		}
		if currentToken != "" || currentToken == " " {
			currentTokenList = append(currentTokenList, currentToken)
		}
		currentToken = ""
	}
}

func Advance() bool {
	if bfr.Scan() {
		lineNumber++
		currenLine = bfr.Text()

		if strings.HasPrefix(currenLine, "//") {
			bfr.Scan()
			lineNumber++
			currenLine = bfr.Text()
		}

		if !isValidParentheses(currenLine) {
			fmt.Println("Invalid parentheses on line ", lineNumber)
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
