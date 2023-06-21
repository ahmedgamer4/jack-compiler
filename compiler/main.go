package main

import (
	"fmt"
	"os"

	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func main() {
	// jacktokenizer.OpenFile(os.Args[1])
	// for jacktokenizer.Advance() {
	// 	c := jacktokenizer.GetCurrentTokens()
	// 	for _, item := range c {
	// 		if !(item == " " || item == "") {
	// 			fmt.Println(item)
	// 		}
	// 	}
	// }
	// jacktokenizer.CloseFile()
	testTokenizer()
}

func testTokenizer() {
	jacktokenizer.OpenFile(os.Args[1])
	for jacktokenizer.Advance() {
		c := jacktokenizer.GetCurrentTokens()
		for _, item := range c {
			if !(item == " " || item == "") {
				opened := "<" + jacktokenizer.GetTokenType(item) + ">"
				closed := "<" + jacktokenizer.GetTokenType(item) + ">"

				fmt.Println(opened + item + closed)
			}
		}
	}
	jacktokenizer.CloseFile()
}
