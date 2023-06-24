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
	input = jacktokenizer.GetCurrentTokensList()
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

func appendClose(tag string) {
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
	handleSyntaxError("Expected", tokenType, str, "got", tag, "on line", jacktokenizer.GetCurrentLineNumber())
}

func identifier() {
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == "identifier" {
		append(tag, currentToken)
		return
	}
	handleSyntaxError("Expected identifier got", tag, "on line", jacktokenizer.GetCurrentLineNumber())
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
	} else {
		identifier() // If you have a custom type this will make sense
	}
}

func compileClass() {
	appendOpen("class")
	eat("class", "keyword")
	eat("className", "identifier")
	eat("{", "symbol")
	compileClassVarDec()
	compileSubroutineDec()
	eat("}", "symbol")
	appendClose("class")
}

func compileClassVarDec() {
	appendOpen("classVarDec")
	if currentToken == "static" {
		eat("static", "keyword")
	} else if currentToken == "field" {
		eat("field", "keyword")
	} else {
		handleSyntaxError("(static | field ) keyword expected on line", jacktokenizer.GetCurrentLineNumber())
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
			handleSyntaxError("Expected symbol , or ; on line", jacktokenizer.GetCurrentLineNumber())
		}
	}
	appendClose("classVarDec")
}

func compileSubroutineDec() {
	appendOpen("subroutine")
	if currentToken == "function" {
		eat("function", "keyword")
	} else if currentToken == "method" {
		eat("method", "keyword")
	} else if currentToken == "constructor" {
		eat("constructor", "keyword")
	} else {
		handleSyntaxError("Expected (function | method | constructor) keyword on line", jacktokenizer.GetCurrentLineNumber())
	}

	handleTypes(true)
	identifier()

	eat("(", "symbol")
	compileParamterList()
	eat(")", "symbol")
	compileSubroutineBody()
	appendClose("subroutine")
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
	appendOpen("subroutineBody")
	eat("{", "symbol")
	compileVarDec()
	compileStatements()
	eat("}", "symbol")
	appendClose("subroutineBody")
}

func compileVarDec() {
	if currentToken == "var" {
		appendOpen("varDec")
		eat("var", "keyword")
		handleTypes(false)

		for i < len(input) {
			identifier()
			if currentToken == "," {
				eat(",", "symbol")
			} else if currentToken == ";" {
				eat(";", "symbol")
				break
			} else {
				handleSyntaxError("Expected symbol , or ; on line", jacktokenizer.GetCurrentLineNumber())
			}
		}

		appendClose("varDec")
	}
}

func compileStatements() {
	appendOpen("statements")
	switch currentToken {
	case "let":
		compileLet()
	case "if":
		compileIf()
	case "while":
		compileWhile()
	case "do":
		compileDo()
	case "return":
		compileReturn()
	case "}":
		return
	default:
		handleSyntaxError("Expected statment (let | if | while | do | return)")
	}
	appendClose("statements")
}

func compileLet() {
	if currentToken == "let" {
		appendOpen("letStatement")
		eat("let", "keyword")
		identifier()
		eat("=", "symbol")
		compileExpression()
		eat(";", "symbol")
		appendClose("letStatement")
	}
}

func compileIf() {
	if currentToken == "if" {
		appendOpen("ifStatement")
		eat("if", "keyword")
		eat("(", "symbol")
		compileExpression()
		eat(")", "symbol")
		eat("{", "symbol")
		compileStatements()
		eat("}", "symbol")
		appendClose("ifStatement")
	}
}

func compileWhile() {
	if currentToken == "while" {
		appendOpen("whileStatement")
		eat("while", "keyword")
		eat("(", "symbol")
		compileExpression()
		eat(")", "symbol")
		eat("{", "symbol")
		compileStatements()
		eat("}", "symbol")
		appendClose("whileStatement")
	}
}

func compileDo() {
	if currentToken == "do" {
		appendOpen("doStatement")
		eat("do", "keyword")
		identifier()
		eat("(", "symbol")
		compileExpressionList()
		eat(")", "symbol")
		eat(";", "symbol")
		appendClose("doStatement")
	}
}

func compileReturn() {
	if currentToken == "return" {
		appendOpen("returnStatement")
		eat("return", "keyword")
		compileExpression()
		eat(";", "symbol")
		appendClose("returnStatement")
	}
}

func compileExpression() {
	appendOpen("expression")
	compileTerm()

	if jacktokenizer.GetTokenType(currentToken) == "symbol" {
		switch currentToken {
		case "+":
			eat("+", "symbol")
		case "-":
			eat("-", "symbol")
		case "/":
			eat("/", "symbol")
		case "*":
			eat("*", "symbol")
		case "&":
			eat("&", "symbol")
		case "|":
			eat("|", "symbol")
		case "<":
			eat("<", "symbol")
		case ">":
			eat(">", "symbol")
		case "=":
			eat("=", "symbol")
		}

		compileTerm()
	}

	appendClose("expression")
}

func compileTerm() {
	appendOpen("term")

	switch jacktokenizer.GetTokenType(currentToken) {
	case "integerConstant":
		append(currentToken, "integerConstant")
		nextToken()
	case "stringConstant":
		append(currentToken, "stringConstant")
		nextToken()
	case "keyword":
		if currentToken == "false" {
			eat("false", "keyword")
		} else if currentToken == "true" {
			eat("true", "keyword")
		} else if currentToken == "null" {
			eat("null", "keyword")
		} else if currentToken == "this" {
			eat("this", "keyword")
		}
	case "identifier":
		identifier()
		if jacktokenizer.GetTokenType(currentToken) == "symbol" {
			if currentToken == "." {
				eat(".", "symbol")
				identifier()
				eat("(", "symbol")
				compileExpressionList()
				eat(")", "symbol")
			} else if currentToken == "[" {
				eat("[", "symbol")
				compileExpression()
				eat("]", "symbol")
			} else if currentToken == "(" {
				eat("(", "symbol")
				compileExpressionList()
				eat(")", "symbol")
			}
		}
	case "symbol":
		if currentToken == "~" {
			eat("~", "symbol")
		} else if currentToken == "-" {
			eat("-", "symbol")
		} else if currentToken == "(" {

		} else {
			handleSyntaxError("Expected ~ | - on line", jacktokenizer.GetCurrentLineNumber())
		}
		compileTerm()
	}

	appendClose("term")
}

func compileExpressionList() {
	appendOpen("expressionList")
	for i < len(input) {
		if currentToken == ")" {
			break
		}
		compileExpression()
		if currentToken == "," {
			eat(",", "symbol")
		}
	}
	appendClose("expressionList")
}
