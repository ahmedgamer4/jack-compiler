package codegenerator

/**
* Simple module that writes VM commands to the output .vm file
**/

import (
	"fmt"
	"os"
)

type Segment string
type Command string

type VMWriter struct {
	File *os.File
}

func (v *VMWriter) Initialize(file *os.File) {
	v.File = file
}

func (v *VMWriter) WritePush(seg Segment, i int) {
	code := fmt.Sprintf(
		`
    push %s %d
    `, seg, i)
	v.File.WriteString(code)
}

func (v *VMWriter) WritePop(seg Segment, i int) {
	code := fmt.Sprintf(
		`
    pop %s %d
    `, seg, i)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteArithmetic(command Command) {
	switch command {
	case "+":
		v.File.WriteString("add\n")
		break
	case "-":
		v.File.WriteString("sub\n")
		break
	case "*":
		v.WriteCall("Math.multiply", 2)
		break
	case "/":
		v.WriteCall("Math.divide", 2)
		break
	case "&":
		v.File.WriteString("and\n")
		break
	case "|":
		v.File.WriteString("or\n")
		break
	case "<":
		v.File.WriteString("lt\n")
		break
	case ">":
		v.File.WriteString("gt\n")
		break
	case "=":
		v.File.WriteString("eq\n")
		break
	case "neg":
		v.File.WriteString("eq\n")
		break
	default:
		println("Unknown operator")
	}
}

func (v *VMWriter) WriteLabel(label string) {
	code := fmt.Sprintf(
		`
    label %s
    `, label)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteGoto(label string) {
	code := fmt.Sprintf(
		`
    goto %s
    `, label)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteIf(label string) {
	code := fmt.Sprintf(
		`
    not
    if-goto %s
    `, label)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteCall(name string, nArgs int) {
	code := fmt.Sprintf(
		`
    call %s %d
    `, name, nArgs)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteFunction(name string, nArgs int) {
	code := fmt.Sprintf("function %s %d\n", name, nArgs)
	v.File.WriteString(code)
}

func (v *VMWriter) WriteReturn() {
	v.File.WriteString("return\n")
}

/**
* Close the output file / stream
* */
func (v *VMWriter) Close() {
	v.File.Close()
}
