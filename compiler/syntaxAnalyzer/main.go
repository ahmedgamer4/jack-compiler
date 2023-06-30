package syntaxanalyzer

import (
	"os"

	compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func HandleFileInput() {
	jacktokenizer.OpenFile(os.Args[1])
	compilationengine.CompileClass()
	jacktokenizer.CloseFile()
	s := compilationengine.GetSytaxTree()

	out := os.Args[1][:len(os.Args[1])-5] + ".test.xml"
	outputFile, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic("Error creating file")
	}

	outputFile.WriteString(s)
	println(s, outputFile)
}
