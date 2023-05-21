package main

import (
	"fmt"
	"strings"
)

var nextLabel = 1

func WriteArithmetic(cRecord CommandRecord) string {
	commands := ""

	switch cRecord.command {
	case "add":
		commands = `
			@SP
			AM=M-1
			D=M
			@SP
			A=M-1
			D=M+D
			M=D
		`
	case "sub":
		commands = `
			@SP
			AM=M-1
			D=M
			@SP
			A=M-1
			D=M-D
			M=D
		`
	case "neg":
		commands = `
			@SP
			A=M-1
			M=-M
		`
	case "not":
		commands = `
			@SP
			A=M-1
			M=!M
		`
	case "and":
		commands = `
			@SP
			AM=M-1
			D=M
			@SP
			A=M-1
			M=D&M
		`
	case "or":
		commands = `
			@SP
			AM=M-1
			D=M
			@SP
			A=M-1
			M=D|M
		`
	case "eq", "gt", "lt":
		commands = fmt.Sprintf(`
			@SP
			AM=M-1
			D=M
			@SP
			A=M-1
			D=M-D
			@eqTrue%d
			D;%s
			@SP
			A=M-1
			M=0
			@eqEnd%d
			0;JMP
			(eqTrue%d)
			@SP
			A=M-1
			M=-1
			(eqEnd%d)
		`, nextLabel, getJumpComparison(cRecord.command), nextLabel, nextLabel, nextLabel)
		nextLabel++
	}

	return "// " + cRecord.command + "\n" + commands
}

func WritePush(cRecord CommandRecord) string {
	commands := ""

	switch cRecord.segment {
	case "constant":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@SP
			A=M
			M=D
			@SP
			M=M+1
		`, cRecord.num)
	case "local", "argument", "this", "that":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@%s
			A=M+D
			D=M
			@SP
			A=M
			M=D
			@SP
			M=M+1
		`, cRecord.num, getSegmentBase(cRecord.segment))

	case "static":
		commands = fmt.Sprintf(`
			@%s.%s
			D=M
			@SP
			A=M
			M=D
			@SP
			M=M+1
		`, strings.Split(GetFilename(), ".")[0], cRecord.num)
	case "temp":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@5
			A=A+D
			D=M
			@SP
			A=M
			M=D
			@SP
			M=M+1
		`, cRecord.num)
	case "pointer":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@3
			A=A+D
			D=M
			@SP
			A=M
			M=D
			@SP
			M=M+1
		`, cRecord.num)
	}

	return "// " + cRecord.command + " " + cRecord.segment + " " + cRecord.num + "\n" + commands
}

func WritePop(cRecord CommandRecord) string {
	commands := ""

	switch cRecord.segment {
	case "local", "argument", "this", "that":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@%s
			D=M+D
			@R13
			M=D
			@SP
			A=M-1
			D=M
			@R13
			A=M
			M=D
		`, cRecord.num, getSegmentBase(cRecord.segment))
	case "static":
		commands = fmt.Sprintf(`
			@SP
			A=M-1
			D=M
			@%s.%s
			M=D
		`, strings.Split(GetFilename(), ".")[0], cRecord.num)
	case "temp":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@5
			D=A+D
			@R13
			M=D
			@SP
			A=M-1
			D=M
			@R13
			A=M
			M=D
		`, cRecord.num)
	case "pointer":
		commands = fmt.Sprintf(`
			@%s
			D=A
			@3
			D=A+D
			@R13
			M=D
			@SP
			A=M-1
			D=M
			@R13
			A=M
			M=D
		`, cRecord.num)
	}

	return "// " + cRecord.command + " " + cRecord.segment + " " + cRecord.num + "\n" + commands
}

func Translate(cRecord CommandRecord) string {
	switch GetCommandType(cRecord.command) {
	case C_PUSH:
		return WritePush(cRecord)
	case C_POP:
		return WritePop(cRecord)
	case C_ARITHMETIC:
		return WriteArithmetic(cRecord)
	default:
		return ""
	}
}

func getSegmentBase(segment string) string {
	switch segment {
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this":
		return "THIS"
	case "that":
		return "THAT"
	default:
		return ""
	}
}

func getJumpComparison(op string) string {
	switch op {
	case "eq":
		return "JEQ"
	case "gt":
		return "JGT"
	case "lt":
		return "JLT"
	default:
		return ""
	}
}
