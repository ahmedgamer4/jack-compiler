package jacktokenizer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	CLASS_METHOD TokenType = iota
	FUNCTION
	CONSTRUCTOR
	INT
	VAR
	BOOLEAN
	CHAR
	VOID
	STATIC
	FIELD
	LET
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

const ()

var (
	lineNumber       int
	pos              int
	bfr              *bufio.Scanner
	currenLine       string
	currentToken     string
	currentTokenList []string

	tokens = map[string]TokenType{
		"keyword": KEYWORD,
	}

	inBlockComment, inString, inCurrly, inBracket bool // These are the parser state

	keywords = map[string]TokenType{
		"class":       CLASS_METHOD,
		"function":    FUNCTION,
		"constructor": CONSTRUCTOR,
		"int":         INT,
		"var":         VAR,
		"boolean":     BOOLEAN,
		"char":        CHAR,
		"void":        VOID,
		"static":      STATIC,
		"field":       FIELD,
		"let":         LET,
		"if":          IF,
		"else":        ELSE,
		"while":       WHILE,
		"return":      RETURN,
		"true":        TRUE,
		"false":       FALSE,
		"null":        NULL,
		"this":        THIS,
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

func GetCurrentTokens() []string {
	return currentTokenList
}

func GetCurrentToken() string {
	return currentToken
}

func GetPos() (int, int) {
	return lineNumber, pos
}

func getTokenType(token string) TokenType {
	if _, ok := keywords[token]; ok {
		return KEYWORD
	} else if _, ok := symbols[token]; ok {
		return SYMBOL
	} else if _, err := strconv.Atoi(token); err == nil {
		return INT_CONST
	} else if strings.HasPrefix(token, "\"") {
		return STRING_CONST
	} else {
		return IDENTIFIER
	}
}

func OpenFile() {
	lineNumber = 1
	pathName := os.Args[1]

	inputInfo, err := os.Stat(pathName)
	handleError(err)

	if inputInfo.IsDir() {
		dir, err := os.ReadDir(pathName)
		handleError(err)

		for _, file := range dir {
			fmt.Println(file)
		}

	} else {
		file, err := os.Open(pathName)
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
						fmt.Println("Missing } on pos ", lineNumber, ", ", pos)
					} else {
						inCurrly = true
					}
				} else if letter == '[' {
					if inBracket {
						fmt.Println("Missing ] on pos ", lineNumber, ", ", pos)
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
						fmt.Println("Missing { on pos ", lineNumber, ", ", pos)
					} else {
						inCurrly = false
					}
				} else if letter == ']' {
					currentToken = string(letter)
					tempToken = ""

					if !inBracket {
						fmt.Println("Missing [ on pos ", lineNumber, ", ", pos)
					} else {
						inBracket = false
					}
				}
			} else if string(letter) == "\n" && inString {
				fmt.Println("Missing \" in position ", lineNumber, ", ", pos)
			}
		}
		currentTokenList = append(currentTokenList, currentToken)
		currentToken = ""
	}
}

func Advance() bool {
	if bfr.Scan() {
		lineNumber++
		currenLine = bfr.Text()
		log.Println(currenLine)

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

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
