package compilationengine

import (
	"fmt"

	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

var (
	i = 0

	syntaxError = false // if true then do not compile but show errors

	input []string

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

	output += "<" + tag + ">" + content + "</" + tag + ">\n"
}

func advance() {
	jacktokenizer.Advance()
	input = jacktokenizer.GetCurrentTokens()
}

func nextToken() {
	if i < len(input) {
		i++
	} else {
		i = 0
		advance()
	}
}

func appnedOpen(tag string) {
	output += "<" + tag + ">\n"
}

func appnedClose(tag string) {
	output += "</" + tag + ">\n"
}

func setCurrentToken(str string) {
	currentToken = str
}

func eat(str string, tokenType string) {
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == tokenType {
		if currentToken == str {
			append(tag, str)
			nextToken()
			return
		}
	}
	fmt.Println("Expected", tokenType, str, "got", tag, "on line", jacktokenizer.GetCurrentLine())
}

func identifier() {
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == "identifier" {
		append(tag, currentToken)
		return
	}
	fmt.Println("Expected identifier got", tag, "on line", jacktokenizer.GetCurrentLine())
}

func handleSyntaxError(message ...any) {
	fmt.Println(message...)
	syntaxError = true
}

func compileClass() {
	eat("class", "keyword")
	eat("className", "identifier")
	eat("{", "symbol")
	// TODO: Add a way to handle class var dec and subroutine
	eat("}", "symbol")
}

func compileClassVarDec() {
	appnedOpen("classVarDec")
	if currentToken == "static" {
		eat("static", "keyword")
	} else if currentToken == "field" {
		eat("field", "keyword")
	} else {
		fmt.Println("(static | field ) keyword expected on line", jacktokenizer.GetCurrentLine())
	}

	if currentToken == "int" {
		eat("int", "keyword")
	} else if currentToken == "char" {
		eat("char", "keyword")
	} else if currentToken == "boolean" {
		eat("boolean", "keyword")
	} else if currentToken == "char" {
		eat("char", "keyword")
	} else {
		identifier() // If you have a custom type this will make sense
	}

	identifier()
	nextToken()
	for i < len(input) {
		if currentToken == "," {
			eat(",", "symbol")
		} else if currentToken == ";" {
			eat(";", "symbol")
		} else {
			handleSyntaxError("expected symbol , or ; on line", jacktokenizer.GetCurrentLine())
		}
	}
	appnedClose("classVarDec")
}
