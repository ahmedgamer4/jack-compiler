package codegenerator

/**
* Simple module that writes VM commands to the output .vm file
**/

import (
	"bytes"
	"fmt"
)

type (
	Segment string
	Command string
)

type VMWriter struct {
	StringBuffer bytes.Buffer
}

func NewWriter() *VMWriter {
	return &VMWriter{}
}

func (v *VMWriter) ResetVmCode() {
	v.StringBuffer.Reset()
}

func (v *VMWriter) WritePush(seg Segment, i int) {
	code := fmt.Sprintf(
		`push %s %d
    `, seg, i)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WritePop(seg Segment, i int) {
	code := fmt.Sprintf(
		`pop %s %d
    `, seg, i)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteArithmetic(command Command) {
	switch command {
	case "+":
		v.StringBuffer.WriteString("add\n")
	case "-":
		v.StringBuffer.WriteString("sub\n")
	case "*":
		v.WriteCall("Math.multiply", 2)
	case "/":
		v.WriteCall("Math.divide", 2)
	case "&":
		v.StringBuffer.WriteString("and\n")
	case "|":
		v.StringBuffer.WriteString("or\n")
	case "<":
		v.StringBuffer.WriteString("lt\n")
	case ">":
		v.StringBuffer.WriteString("gt\n")
	case "=":
		v.StringBuffer.WriteString("eq\n")
	case "~":
		v.StringBuffer.WriteString("not\n")
	case "neg":
		v.StringBuffer.WriteString("neg\n")
	default:
		println("Unknown operator")
	}
}

func (v *VMWriter) WriteLabel(label string) {
	code := fmt.Sprintf(
		`label %s
    `, label)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteGoto(label string) {
	code := fmt.Sprintf(
		`goto %s
    `, label)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteIf(label string) {
	code := fmt.Sprintf(
		`if-goto %s
    `, label)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteCall(name string, nArgs int) {
	code := fmt.Sprintf(
		`call %s %d
    `, name, nArgs)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteFunction(name string, nArgs int) {
	code := fmt.Sprintf("function %s %d\n", name, nArgs)
	v.StringBuffer.WriteString(code)
}

func (v *VMWriter) WriteReturn() {
	v.StringBuffer.WriteString("return\n")
}

func (v *VMWriter) GetVMCode() string {
	return v.StringBuffer.String()
}
