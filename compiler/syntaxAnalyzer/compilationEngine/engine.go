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

func appendOpen(tag string) {
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
	handleSyntaxError("Expected", tokenType, str, "got", tag, "on line", jacktokenizer.GetCurrentLine())
}

func identifier() {
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == "identifier" {
		append(tag, currentToken)
		return
	}
	handleSyntaxError("Expected identifier got", tag, "on line", jacktokenizer.GetCurrentLine())
}

func handleSyntaxError(message ...any) {
	fmt.Println(message...)
	syntaxError = true
}

func handleTypes(isFunction bool) {
	if isFunction {
		if currentToken == "void" {
			eat("void", "keyword")
		}
	} else if currentToken == "int" {
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
}

func compileClass() {
	appendOpen("class")
	eat("class", "keyword")
	eat("className", "identifier")
	eat("{", "symbol")
	// TODO: Add a way to handle class var dec and subroutine
	compileClassVarDec()
	eat("}", "symbol")
	appnedClose("class")
}

func compileClassVarDec() {
	appendOpen("classVarDec")
	if currentToken == "static" {
		eat("static", "keyword")
	} else if currentToken == "field" {
		eat("field", "keyword")
	} else {
		handleSyntaxError("(static | field ) keyword expected on line", jacktokenizer.GetCurrentLine())
	}

	handleTypes(false)
	for i < len(input) {
		identifier()
		if currentToken == "," {
			eat(",", "symbol")
		} else if currentToken == ";" {
			eat(";", "symbol")
			break
		} else {
			handleSyntaxError("Expected symbol , or ; on line", jacktokenizer.GetCurrentLine())
		}
	}
	appnedClose("classVarDec")
}

func compileSubroutine() {
	appendOpen("subroutine")
	if currentToken == "function" {
		eat("function", "keyword")
	} else if currentToken == "method" {
		eat("method", "keyword")
	} else if currentToken == "constructor" {
		eat("constructor", "keyword")
	} else {
		handleSyntaxError("Expected (function | method | constructor) keyword on line", jacktokenizer.GetCurrentLine())
	}

	handleTypes(true)
	identifier()

	eat("(", "symbol")
	compileParamterList()
	eat(")", "symbol")
	eat("{", "symbol")
	compileSubroutineBody()
	eat("}", "symbol")
	appnedClose("subroutine")
}

func compileParamterList() {
	for currentToken != ")" {
		handleTypes(false)
		identifier()
		if currentToken == "," {
			eat(",", "symbol")
		}
	}
}

func compileSubroutineBody() {
	return
}

func compileVarDec() {
	appendOpen("varDec")
	if currentToken == "var" {
		eat("var", "keyword")
	} else {
		handleSyntaxError("Expected keyword var on line", jacktokenizer.GetCurrentLine())
	}

	handleTypes(false)

	for i < len(input) {
		identifier()
		if currentToken == "," {
			eat(",", "symbol")
		} else if currentToken == ";" {
			eat(";", "symbol")
			break
		} else {
			handleSyntaxError("Expected symbol , or ; on line", jacktokenizer.GetCurrentLine())
		}
	}

	appnedClose("varDec")
}
