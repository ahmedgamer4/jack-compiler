package main

import (
	"fmt"
	"time"

	syntaxanalyzer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer"
	// compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
)

func main() {
	now := time.Now()
	syntaxanalyzer.StartParsing()
	// println(compilationengine.GetVMCode())
	fmt.Println(time.Since(now))
}
