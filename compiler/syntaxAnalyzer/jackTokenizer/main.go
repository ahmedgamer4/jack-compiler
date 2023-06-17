package jacktokenizer

import (
	"bufio"
	"fmt"
	"os"
)

type TokenType int
type KeywordType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	CLASS_METHOD KeywordType = iota
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

var (
	lineNumber int
	bfr        *bufio.Scanner
)

func OpenFile() {
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

		for advance() {
			fmt.Println(bfr.Text())
		}
	}
}

func advance() bool {
	lineNumber++
	return bfr.Scan()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
