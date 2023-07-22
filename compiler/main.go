package main

import (
	"fmt"
	"time"

	syntaxanalyzer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer"
)

func main() {
	now := time.Now()
	syntaxanalyzer.StartParsing()
	fmt.Println(time.Since(now))
}
