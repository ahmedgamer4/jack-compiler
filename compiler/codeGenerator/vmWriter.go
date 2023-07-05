package codegenerator

import (
	"fmt"
	"os"
)

/**
* Simple module that writes VM commands to the output .vm file
* */

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
