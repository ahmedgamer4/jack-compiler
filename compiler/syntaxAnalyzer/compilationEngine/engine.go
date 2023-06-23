package compilationengine

import (
	"fmt"

	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

var (
	currentToken = ""
	output       = ""
)

// TODO: Compelete this function: it should start the whole engine
func Run() {
	return
}

func append(tag string, content string) {
	if tag == "&" {
		tag = "&amp;"
	} else if tag == "<" {
		tag = "&lt;"
	} else if tag == ">" {
		tag = "&gt;"
	}

	content += "<" + tag + ">" + content + "</" + tag + ">\n"
}

func eat(str string, tokenType string) {
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == tokenType {
		if currentToken == str {
			append(tag, str)
			return
		}
	}
	fmt.Println("Expected", tokenType, str, "got", tag)
}

func compileClass() {
	eat("class", "keyword")
	eat("className", "identifier")
	eat("{", "symbol")
	// TODO: Add a way to handle class var dec and subroutine
	eat("}", "symbol")
}
