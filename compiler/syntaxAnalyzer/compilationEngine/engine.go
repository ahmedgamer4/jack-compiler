package compilationengine

import (
	"fmt"
	"strings"

	codegenerator "github.com/ahmedgamer4/jack-compiler/compiler/codeGenerator"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"
)

var (
	i = 0

	syntaxError = false // if true then do not compile but show errors

	input []string

	currentToken = ""
	syntaxTree   = ""

	symbolTable = codegenerator.SymbolTable{}
	writer      = codegenerator.VMWriter{}
)

func IsSyntaxError() bool {
	return syntaxError
}

func appendTag(tag string, content string) {
	if content == "&" {
		content = "&amp;"
	} else if content == "<" {
		content = "&lt;"
	} else if content == ">" {
		content = "&gt;"
	}

	syntaxTree += "<" + tag + "> " + content + "" + " </" + tag + ">\n"
}

func GetSytaxTree() string {
	return syntaxTree
}

/**
* advance in terms of line not the token
* in order to know the current line
* */
func advance() bool {
	advance := jacktokenizer.Advance()
	input = jacktokenizer.GetCurrentTokensList()
	for len(input) == 0 {
		jacktokenizer.Advance()
		input = jacktokenizer.GetCurrentTokensList()
	}
	currentToken = input[i]
	return advance
}

func nextToken() bool {
	if i < len(input)-1 {
		i++
		currentToken = strings.TrimSpace(input[i])
		return true
	}
	i = 0
	adv := advance()
	currentToken = strings.TrimSpace(input[i])
	return adv
}

func appendOpen(tag string) {
	syntaxTree += "<" + tag + ">\n"
}

func appendClose(tag string) {
	syntaxTree += "</" + tag + ">\n"
}

func isTokenEmpty() {
	for currentToken == "" || currentToken == " " {
		nextToken()
	}
}

func eat(str string, tokenType string) {
	isTokenEmpty()
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == tokenType {
		if currentToken == str {
			appendTag(tag, currentToken)
			nextToken()
			isTokenEmpty()
			return
		}
	}
	handleSyntaxError("Expected", tokenType, str, "got", tag, "on line", jacktokenizer.GetCurrentLineNumber(), currentToken)
}

func identifier() {
	isTokenEmpty()
	tag := jacktokenizer.GetTokenType(currentToken)
	if tag == "identifier" {
		appendTag(tag, currentToken)
		nextToken()
		isTokenEmpty()
		return
	}
	handleSyntaxError("Expected identifier got", tag, "on line", jacktokenizer.GetCurrentLineNumber(), currentToken, input)
}

func handleSyntaxError(message ...any) {
	fmt.Println(message, jacktokenizer.GetCurrentLineNumber())
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

func CompileClass() {
	advance()
	isTokenEmpty()
	if currentToken != "class" {
		handleSyntaxError("Expected keyword class on line", jacktokenizer.GetCurrentLineNumber())
	}
	appendOpen("class")
	eat("class", "keyword")
	identifier()
	eat("{", "symbol")

	for currentToken == "static" || currentToken == "field" {
		compileClassVarDec()
	}

	for currentToken == "function" || currentToken == "method" || currentToken == "constructor" {
		compileSubroutineDec()
	}

	eat("}", "symbol")
	appendClose("class")
}

func compileClassVarDec() {
	var currentKind codegenerator.FieldType
	var currentType string

	appendOpen("classVarDec")
	if currentToken == "static" {
		currentKind = codegenerator.FieldType(currentToken)
		eat("static", "keyword")
	} else if currentToken == "field" {
		currentKind = codegenerator.FieldType(currentToken)
		eat("field", "keyword")
	} else {
		handleSyntaxError("(static | field ) keyword expected on line", jacktokenizer.GetCurrentLineNumber())
	}

	currentType = currentToken
	handleTypes(false)

	for i < len(input) {
		symbolTable.Define(currentToken, currentType, currentKind)
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
	fmt.Println(symbolTable.ClassSymbolTable, jacktokenizer.GetCurrentLineNumber())
	appendClose("classVarDec")
}

func compileSubroutineDec() {
	appendOpen("subroutineDec")
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
	println(currentToken)
	identifier()
	println(currentToken)

	eat("(", "symbol")
	compileParameterList()
	eat(")", "symbol")
	compileSubroutineBody()
	appendClose("subroutineDec")
}

func compileParameterList() {
	appendOpen("parameterList")

	for currentToken != ")" {
		handleTypes(false)
		identifier()
		if currentToken == "," {
			eat(",", "symbol")
		}
	}
	appendClose("parameterList")
}

func compileSubroutineBody() {
	appendOpen("subroutineBody")
	eat("{", "symbol")
	for currentToken == "var" {
		compileVarDec()
	}
	for currentToken == "if" || currentToken == "let" || currentToken == "while" || currentToken == "return" || currentToken == "do" {
		compileStatements()
	}
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
				handleSyntaxError("Expected symbol , or ; on line", jacktokenizer.GetCurrentLineNumber(), "got", currentToken)
			}
		}

		appendClose("varDec")
	}
}

func compileStatements() {
	appendOpen("statements")
	for currentToken == "if" || currentToken == "let" || currentToken == "while" || currentToken == "return" || currentToken == "do" {
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
			break
		default:
			handleSyntaxError("Expected statment (let | if | while | do | return) got", currentToken)
		}
	}
	appendClose("statements")
}

func compileLet() {
	if currentToken == "let" {
		appendOpen("letStatement")
		eat("let", "keyword")
		identifier()
		for currentToken == "[" {
			eat("[", "symbol")
			compileExpression()
			eat("]", "symbol")
		}
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
		if currentToken == "else" {
			eat("else", "keyword")
			eat("{", "symbol")
			compileStatements()
			eat("}", "symbol")
		}
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
		if currentToken == "." {
			eat(".", "symbol")
			identifier()
		}
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
		if currentToken == ";" {
			eat(";", "symbol")
		} else {
			compileExpression()
			eat(";", "symbol")
		}
		appendClose("returnStatement")
	}
}

func compileExpression() {
	appendOpen("expression")
	compileTerm()

	// invalidSymbols := map[string]int{
	// 	"{": 0,
	// 	"}": 0,
	// 	"(": 0,
	// 	")": 0,
	// 	"[": 0,
	// 	"]": 0,
	// 	".": 0,
	// 	",": 0,
	// 	";": 0,
	// }

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
		default:
			appendClose("expression")
			return
		}

		compileTerm()
	}

	appendClose("expression")
}

func compileTerm() {
	if currentToken == ";" || currentToken == " " {
		return
	}
	appendOpen("term")

	switch jacktokenizer.GetTokenType(currentToken) {
	case "integerConstant":
		appendTag("integerConstant", currentToken)
		nextToken()
	case "stringConstant":
		appendTag("stringConstant", currentToken)
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
		} else if currentToken == "return" {
			eat("return", "keyword")
		}
	case "identifier":
		identifier()
		if jacktokenizer.GetTokenType(currentToken) == "symbol" || currentToken != ";" {
			if currentToken == "." {
				eat(".", "symbol")
				identifier()
				eat("(", "symbol")
				compileExpressionList()
				eat(")", "symbol")
			} else if currentToken == "[" {
				for currentToken == "[" {
					eat("[", "symbol")
					compileExpression()
					eat("]", "symbol")
				}
			} else if currentToken == "(" {
				eat("(", "symbol")
				compileExpressionList()
				eat(")", "symbol")
			}
		}
		// TODO: Fix this function
	case "symbol":
		if currentToken == "~" {
			eat("~", "symbol")
			compileTerm()
		} else if currentToken == "-" {
			eat("-", "symbol")
			compileTerm()
		} else if currentToken == "(" {
			eat("(", "symbol")
			compileExpression()
			eat(")", "symbol")
		} else {
			handleSyntaxError("Expected ~ | - on line", jacktokenizer.GetCurrentLineNumber())
		}
	default:
		handleSyntaxError("symbol | identifier | string | integer expected")
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
