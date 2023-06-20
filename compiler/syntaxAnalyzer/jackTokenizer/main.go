package jacktokenizer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	lineNumber   int
	pos          int
	bfr          *bufio.Scanner
	currenLine   string
	currentToken string

	tokens = map[string]TokenType{
		"keyword": KEYWORD,
	}

	blockComment, inString, isCurrly, isBracket bool // These are the parser state

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

func ReturnCurrentToken() string {
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

		defer file.Close()

		bfr = bufio.NewScanner(file)
		bfr.Split(bufio.ScanLines)
	}
}

func handleTokens(currenLine string) {
	tempToken := ""

	for i, letter := range currenLine {
		tempToken += string(letter)
		pos = i
		if letter == ' ' {
			continue
		} else {
			if _, ok := keywords[tempToken]; ok {
				currentToken = tempToken
				tempToken = ""
			} else if letter == '{' {
				if isCurrly {
					fmt.Println("Missing } on pos ", lineNumber, ", ", pos)
				} else {
					isCurrly = true
				}
			} else if letter == '[' {
				if isBracket {
					fmt.Println("Missing ] on pos ", lineNumber, ", ", pos)
				} else {
					isBracket = true
				}
			} else if _, ok := symbols[string(letter)]; ok {
				currentToken = string(letter)
				tempToken = ""
			} else if letter == '}' {
				currentToken = string(letter)

				if !isCurrly {
					fmt.Println("Missing { on pos ", lineNumber, ", ", pos)
				} else {
					isCurrly = false
				}
			} else if letter == ']' {
				currentToken = string(letter)

				if !isBracket {
					fmt.Println("Missing [ on pos ", lineNumber, ", ", pos)
				} else {
					isBracket = false
				}
			} else if letter == '*' && currenLine[i+1] == '/' {
				if !blockComment {
					fmt.Println("Missing /* on position ", lineNumber, ", ", pos)
				} else {
					blockComment = false
				}
			} else if letter == '/' && currenLine[i+1] == '*' {
				if blockComment {
					fmt.Println("Missing */ on position ", lineNumber, ", ", pos)
				} else {
					blockComment = true
				}
			} else if letter == '"' {
				if inString {
					inString = false
					currentToken = "\"" + tempToken + "\""
				} else {
					inString = true
				}
			} else if string(letter) == "\n" && inString {
				fmt.Println("Missing \" in pos ", lineNumber, ", ", pos)
			} else {
				tempToken += string(letter)
			}
		}
	}
}

func advance() {

	if bfr.Scan() {
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
		return
	}
	lineNumber++
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
