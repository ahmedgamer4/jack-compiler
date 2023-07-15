package compilationengine

import (
	"fmt"
	"strconv"
	"strings"

	codegenerator "github.com/ahmedgamer4/jack-compiler/compiler/codeGenerator"
	jacktokenizer "github.com/ahmedgamer4/jack-compiler/compiler/syntaxAnalyzer/jackTokenizer"

	"github.com/google/uuid"
)

var (
	i = 0

	syntaxError = false // if true then do not compile but show errors

	input []string

	currentToken = ""
	syntaxTree   = ""

	symbolTable = codegenerator.NewSTable()
	writer      = codegenerator.NewWriter()

	currentCompilingFile = ""
)

func SetCurrentCFile(filename string) {
	currentCompilingFile = filename
}

func IsSyntaxError() bool {
	return syntaxError
}

func GetVMCode() string {
	return writer.GetVMCode()
}

func appendTag(tag string, content string) {
	switch tag {
	case "&":
		content = "&amp;"
	case "<":
		content = "&lt;"
	case ">":
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

func handleSyntaxError(message ...interface{}) {
	fmt.Println(currentCompilingFile, ": ", message, jacktokenizer.GetCurrentLineNumber(), jacktokenizer.GetCurrentTokensList())
	syntaxError = true
}

func handleTypes() {
	switch currentToken {
	case "void":
		eat("void", "keyword")
	case "int":
		eat("int", "keyword")
	case "char":
		eat("char", "keyword")
	case "boolean":
		eat("boolean", "keyword")
	default:
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
	fmt.Println(symbolTable)
}

func compileClassVarDec() {
	var currentKind codegenerator.FieldType
	var currentType string

	appendOpen("classVarDec")
	switch currentToken {
	case "static":
		currentKind = codegenerator.FieldType(currentToken)
		eat("static", "keyword")
	case "field":
		currentKind = codegenerator.FieldType(currentToken)
		eat("field", "keyword")
	default:
		handleSyntaxError("(static | field ) keyword expected on line", jacktokenizer.GetCurrentLineNumber())
	}

	currentType = currentToken
	handleTypes()

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
	appendClose("classVarDec")
}

func compileSubroutineDec() {
	appendOpen("subroutineDec")
	switch currentToken {
	case "function":
		eat("function", "keyword")
	case "method":
		eat("method", "keyword")
	case "constructor":
		eat("constructor", "keyword")
	default:
		handleSyntaxError("Expected (function | method | constructor) keyword on line", jacktokenizer.GetCurrentLineNumber())
	}

	handleTypes()
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
		currentType := currentToken
		handleTypes()
		currentParam := currentToken
		identifier()
		symbolTable.Define(currentParam, currentType, "arg")
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
	compileStatements()
	eat("}", "symbol")
	appendClose("subroutineBody")
}

func compileVarDec() {
	if currentToken == "var" {
		appendOpen("varDec")
		eat("var", "keyword")
		handleTypes()

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
			return
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

		currentVarKind := codegenerator.Segment(symbolTable.KindOf(currentToken))
		currentVarIdx := symbolTable.IndexOf(currentToken)

		identifier()
		for currentToken == "[" {
			eat("[", "symbol")
			compileExpression()
			eat("]", "symbol")
		}
		eat("=", "symbol")
		// You should compile the expression after the "=" sign then pop the result to the current variable
		compileExpression()

		writer.WritePop(currentVarKind, currentVarIdx)

		eat(";", "symbol")
		appendClose("letStatement")
	}
}

// TODO: Do not forget to finished this function
func compileIf() {
	l1 := uuid.New().String()
	l2 := uuid.New().String()

	if currentToken == "if" {
		appendOpen("ifStatement")
		eat("if", "keyword")
		eat("(", "symbol")
		compileExpression()
		writer.WriteArithmetic(codegenerator.Command("~"))
		eat(")", "symbol")
		writer.WriteIf(l1)
		eat("{", "symbol")
		compileStatements()

		writer.WriteGoto(l2)

		eat("}", "symbol")

		writer.WriteLabel(l1)
		if currentToken == "else" {
			eat("else", "keyword")
			eat("{", "symbol")
			compileStatements()
			eat("}", "symbol")
			writer.WriteLabel(l2)
		}
		appendClose("ifStatement")
	}
}

func compileWhile() {
	l1 := uuid.New().String()
	l2 := uuid.New().String()

	if currentToken == "while" {
		appendOpen("whileStatement")
		eat("while", "keyword")
		writer.WriteLabel(l1)
		eat("(", "symbol")
		compileExpression()

		writer.WriteArithmetic(codegenerator.Command("~"))
		writer.WriteIf(l2)
		eat(")", "symbol")
		eat("{", "symbol")
		compileStatements()
		eat("}", "symbol")

		writer.WriteGoto(l1)
		writer.WriteLabel(l2)
		appendClose("whileStatement")
	}
}

func compileDo() {
	var fn string
	if currentToken == "do" {
		appendOpen("doStatement")
		eat("do", "keyword")
		fn += currentToken
		identifier()
		if currentToken == "." {
			eat(".", "symbol")
			fn += "." + currentToken
			identifier()
		}
		eat("(", "symbol")
		nArgs := compileExpressionList()
		eat(")", "symbol")
		eat(";", "symbol")
		writer.WriteCall(fn, nArgs)
		writer.WritePop("temp", 0)
		appendClose("doStatement")
	}
}

func compileReturn() {
	if currentToken == "return" {
		appendOpen("returnStatement")
		eat("return", "keyword")
		isTokenEmpty()
		if currentToken == ";" {
			writer.WritePush("constant", 0)
		}
		compileExpression()
		writer.WriteReturn()
		eat(";", "symbol")
		appendClose("returnStatement")
	}
}

func compileExpression() {
	isTokenEmpty()

	if currentToken == ";" {
		return
	}
	appendOpen("expression")
	fstTerm, _ := strconv.Atoi(currentToken)

	writer.WritePush("constant", fstTerm)
	compileTerm()

	for jacktokenizer.GetTokenType(currentToken) == "symbol" {
		op := currentToken
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
			fmt.Println("appendClose", currentToken, jacktokenizer.GetCurrentTokensList())
			appendClose("expression")
			return
		}

		isTokenEmpty()
		secTerm, err := strconv.Atoi(currentToken)
		if err != nil {
			handleSyntaxError("Cannot add", currentToken, "to int")
		}
		writer.WritePush("constant", secTerm)
		writer.WriteArithmetic(codegenerator.Command(op))
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
		isTokenEmpty()
		i, err := strconv.Atoi(currentToken)
		if err != nil {
			panic("Error converting a string into integer")
		}
		writer.WritePush(codegenerator.Segment("constant"), i)
		appendTag("integerConstant", currentToken)
		nextToken()
	case "stringConstant":
		isTokenEmpty()
		appendTag("stringConstant", currentToken)
		nextToken()
	case "keyword":
		switch currentToken {
		case "false":
			writer.WritePush(codegenerator.Segment("constant"), 0)
			eat("false", "keyword")
		case "true":
			writer.WritePush(codegenerator.Segment("constant"), 1)
			writer.WriteArithmetic("neg")
			eat("true", "keyword")
		case "null":
			writer.WritePush(codegenerator.Segment("constant"), 0)
			eat("null", "keyword")
		case "this":
			writer.WritePush(codegenerator.Segment("pointer"), 0)
			eat("this", "keyword")
		case "return":
			writer.WriteReturn()
			eat("return", "keyword")
		default:
			return
		}
	case "identifier":
		idn := currentToken
		identifier()
		isTokenEmpty()
		if jacktokenizer.GetTokenType(currentToken) == "symbol" || currentToken != ";" {
			switch currentToken {
			case ".":
				idn += currentToken
				eat(".", "symbol")
				identifier()
				eat("(", "symbol")
				nArgs := compileExpressionList()
				eat(")", "symbol")
				writer.WriteCall(idn, nArgs)
			case "[":
				for currentToken == "[" {
					eat("[", "symbol")
					compileExpression()
					eat("]", "symbol")
				}
			case "(":
				eat("(", "symbol")
				nArgs := compileExpressionList()
				eat(")", "symbol")
				writer.WriteCall(idn, nArgs)
			default:
				break
			}
		}
		// TODO: Fix this function
	case "symbol":
		switch currentToken {
		case "~":
			eat("~", "symbol")
			compileTerm()
			writer.WriteArithmetic(codegenerator.Command("~"))
		case "-":
			eat("-", "symbol")
			compileTerm()
			writer.WriteArithmetic(codegenerator.Command("neg"))
		case "(":
			eat("(", "symbol")
			compileExpression()
			eat(")", "symbol")
		default:
			handleSyntaxError("Expected ~ | - on line", jacktokenizer.GetCurrentLineNumber())
		}
	default:
		handleSyntaxError("symbol | identifier | string | integer expected")
	}

	appendClose("term")
}

/**
* It should return an int representing the number of arguments passed
* */
func compileExpressionList() int {
	nArgs := 0
	appendOpen("expressionList")
	for i < len(input) {
		if currentToken == ")" {
			break
		}
		compileExpression()
		nArgs++
		if currentToken == "," {
			eat(",", "symbol")
		}
	}
	appendClose("expressionList")
	return nArgs
}
