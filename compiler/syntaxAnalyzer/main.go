package syntaxanalyzer

import (
	"os"
	"path"

	compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func StartParsing() {
	println("Compiling...")

	inputInfo, err := os.Stat(os.Args[1])
	handleError(err)
	if inputInfo.IsDir() {
		files, err := os.ReadDir(os.Args[1])
		handleError(err)

		for _, file := range files {
			var filePath string

			compilationengine.SetCurrentCFile(file.Name())
			jacktokenizer.SetCurrentCFile(file.Name())

			if os.Args[1][len(os.Args[1])-1] == '/' {
				filePath = os.Args[1] + file.Name()
			} else {
				filePath = os.Args[1] + "/" + file.Name()
			}

			if path.Ext(file.Name()) == ".jack" {
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
	vmCode := compilationengine.GetVMCode()

	out := file[:len(file)-5] + ".test.vm"
	println(out)

	outputFile, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	handleError(err)

	if !compilationengine.IsSyntaxError() {
		outputFile.WriteString(vmCode)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
