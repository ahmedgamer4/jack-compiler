package main

import (
	syntaxanalyzer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer"
	// compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
)

func main() {
	syntaxanalyzer.StartParsing()
	// println(compilationengine.GetVMCode())
}
