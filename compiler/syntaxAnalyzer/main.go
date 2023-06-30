package syntaxanalyzer

import (
	"os"
	"path"

	compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func StartParsing() {
	inputInfo, err := os.Stat(os.Args[1])
	handleError(err)
	if inputInfo.IsDir() {
		files, err := os.ReadDir(os.Args[1])
		handleError(err)

		for _, file := range files {
			var filePath string
			if os.Args[1][len(os.Args[1])-1] == '/' {
				filePath = os.Args[1] + file.Name()
			}
			filePath = os.Args[1] + "/" + file.Name()
			if path.Ext(file.Name()) == ".jack" {
				println(filePath)
				jacktokenizer.OpenFile(filePath)
				createFile(filePath)
			}
		}
	} else {
		if path.Ext(os.Args[1]) != ".jack" {
			panic(".jack files only")
		}
		jacktokenizer.OpenFile(os.Args[1])
		createFile(os.Args[1])
	}
}

func createFile(file string) {
	compilationengine.CompileClass()
	jacktokenizer.CloseFile()
	s := compilationengine.GetSytaxTree()

	out := file[:len(file)-5] + ".test.xml"
	println(out)
	outputFile, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	handleError(err)

	outputFile.WriteString(s)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
