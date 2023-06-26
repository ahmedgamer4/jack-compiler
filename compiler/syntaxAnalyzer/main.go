package syntaxanalyzer

import (
	"os"

	compilationengine "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/compilationEngine"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

func HandleFileInput() {
	jacktokenizer.OpenFile(os.Args[1])
	compilationengine.CompileClass()
	// for jacktokenizer.Advance() {
	// 	fmt.Println(jacktokenizer.GetCurrentTokensList())
	// }
	jacktokenizer.CloseFile()
	s := compilationengine.GetSytaxTree()
	println(s)
}
