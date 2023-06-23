package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	start := time.Now()
	testTokenizer()
	log.Println(time.Since(start))
}

func testTokenizer() {
	jacktokenizer.OpenFile(os.Args[1])
	for jacktokenizer.Advance() {
		c := jacktokenizer.GetCurrentTokens()
		for _, item := range c {
			item = strings.TrimSpace(item)
			if !(item == " " || item == "") {
				if item == "&" {
					item = "&amp;"
				} else if item == "<" {
					item = "&lt;"
				} else if item == ">" {
					item = "&gt;"
				}

				opened := "<" + jacktokenizer.GetTokenType(item) + ">"
				closed := "<" + jacktokenizer.GetTokenType(item) + ">"

				fmt.Println(opened + item + closed)
			}
		}
	}
	jacktokenizer.CloseFile()
}
