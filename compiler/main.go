package main

import (
	"fmt"

	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func main() {
	jacktokenizer.OpenFile()
	for jacktokenizer.Advance() {
		c := jacktokenizer.GetCurrentTokens()
		for _, item := range c {
			if !(item == " " || item == "") {
				fmt.Println(item)
			}
		}
	}
}
