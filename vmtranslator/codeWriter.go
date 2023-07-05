package main

import (
	"fmt"
	"path"
	"strings"
)

var (
	nextLabel = 1
	nextFun   = 1
)

func Bootstrap() string {
	bootstrapRecord := CommandRecord{
		command: "call",
		arg0:    "Sys.init",
		arg1:    "0",
	}

	setupCode := `
  @256
  D=A
  @SP
  M=D
  ` + writeCall(bootstrapRecord)
	return setupCode
}

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
    M=M+D
    `
	case "sub":
		commands = `
    @SP
    AM=M-1
    D=M
    @SP
    A=M-1
    M=M-D
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
      M=-1
      @eqTrue%d
      D;%s
      @SP
      A=M-1
      M=0
      (eqTrue%d)
      `, nextLabel, getJumpComparison(cRecord.command), nextLabel)
		nextLabel++
	}

	return "// " + cRecord.command + "\n" + commands
}

func WritePush(cRecord CommandRecord) string {
	commands := ""

	switch cRecord.arg0 {
	case "constant":
		commands = fmt.Sprintf(`
      @%s
      D=A
      @SP
      A=M
      M=D
      @SP
      M=M+1
      `, cRecord.arg1)
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
      `, cRecord.arg1, getSegmentBase(cRecord.arg0))

	case "static":
		commands = fmt.Sprintf(`
      @%s.%s
      D=M
      @SP
      A=M
      M=D
      @SP
      M=M+1
      `, strings.Split(path.Base(GetInArg()), ".")[0], cRecord.arg1)
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
      `, cRecord.arg1)
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
      `, cRecord.arg1)
	}

	return "// " + cRecord.command + " " + cRecord.arg0 + " " + cRecord.arg1 + "\n" + commands
}

func WritePop(cRecord CommandRecord) string {
	commands := ""

	switch cRecord.arg0 {
	case "local", "argument", "this", "that":
		commands = fmt.Sprintf(`
      @%s
      D=A
      @%s
      D=M+D
      @R13
      M=D
      @SP
      AM=M-1
      D=M
      @R13
      A=M
      M=D
      `, cRecord.arg1, getSegmentBase(cRecord.arg0))
	case "static":
		commands = fmt.Sprintf(`
      @SP
      AM=M-1
      D=M
      @%s.%s
      M=D
      `, strings.Split(path.Base(GetInArg()), ".")[0], cRecord.arg1)
	case "temp":
		commands = fmt.Sprintf(`
      @%s
      D=A
      @5
      D=A+D
      @R13
      M=D
      @SP
      AM=M-1
      D=M
      @R13
      A=M
      M=D
      `, cRecord.arg1)
	case "pointer":
		commands = fmt.Sprintf(`
      @%s
      D=A
      @3
      D=A+D
      @R13
      M=D
      @SP
      AM=M-1
      D=M
      @R13
      A=M
      M=D
      `, cRecord.arg1)
	}

	return "// " + cRecord.command + " " + cRecord.arg0 + " " + cRecord.arg1 + "\n" + commands
}

func writeLabel(cRecord CommandRecord) string {
	return fmt.Sprintf(`
    (%s)
    `, cRecord.arg1)
}

func writeGoto(cRecord CommandRecord) string {
	return fmt.Sprintf(`
    @%s
    0;JMP
    `, cRecord.arg1)
}

func writeIfGoto(cRecord CommandRecord) string {
	fmt.Println(cRecord)
	return fmt.Sprintf(`
    @SP
    AM=M-1
    D=M
    @%s
    D;JNE
    `, cRecord.arg1)
}

func writeCall(cRecord CommandRecord) string {
	funName := cRecord.arg0
	nArgs := cRecord.arg1

	res := fmt.Sprintf(`
    @%sret%d
    D=A
    @SP
    A=M
    M=D
    @SP
    M=M+1
    %s
    %s
    %s
    %s
    @SP
    D=M
    @5
    D=D-A
    @%s
    D=D-A
    @ARG
    M=D
    @SP
    D=M
    @LCL
    M=D
    @%s
    0;JMP
    (%sret%d)
    `, funName, nextFun, pushSegIntoStack("LCL"), pushSegIntoStack("ARG"), pushSegIntoStack("THIS"), pushSegIntoStack("THAT"), nArgs, funName, funName, nextFun)

	nextFun++

	return res
}

func writeFun(cRecord CommandRecord) string {
	funName := cRecord.arg0
	nArgs := cRecord.arg1

	res := fmt.Sprintf(`
    (%s)
    `, funName)

	for range nArgs {
		res += pushToStack("0")
	}

	return res
}

func writeReturn(cRecord CommandRecord) string {
	return fmt.Sprintf(`
    @LCL
    D=M
    @R13
    M=D
    @5
    D=A
    @R13
    A=M-D
    D=M
    @retAddr
    M=D
    %s
    @ARG
    D=M+1
    @SP
    M=D
    @R13
    AM=M-1
    D=M
    @THAT
    M=D
    @R13
    AM=M-1
    D=M
    @THIS
    M=D
    @R13
    AM=M-1
    D=M
    @ARG
    M=D
    @R13
    AM=M-1
    D=M
    @LCL
    M=D
    @retAddr
    A=M
    0;JMP
    `, popSegFromStack("ARG"))
}

func pushToStack(item string) string {
	return fmt.Sprintf(`
    @%s
    D=A
    @SP
    A=M
    M=D
    @SP
    M=M+1
    `, item)
}

func pushSegIntoStack(seg string) string {
	return fmt.Sprintf(`
    @%s
    D=M
    @SP
    A=M
    M=D
    @SP
    M=M+1
    `, seg)
}

func popSegFromStack(seg string) string {
	return fmt.Sprintf(`
    @SP
    AM=M-1
    D=M
    @%s
    A=M
    M=D
    `, seg)
}

func Translate(cRecord CommandRecord) string {
	switch GetCommandType(cRecord.command) {
	case C_PUSH:
		return WritePush(cRecord)
	case C_POP:
		return WritePop(cRecord)
	case C_ARITHMETIC:
		return WriteArithmetic(cRecord)
	case C_LABEL:
		return writeLabel(cRecord)
	case C_GOTO:
		return writeGoto(cRecord)
	case C_IF_GOTO:
		return writeIfGoto(cRecord)
	case C_CALL:
		return writeCall(cRecord)
	case C_FUNCTION:
		return writeFun(cRecord)
	case C_RETURN:
		return writeReturn(cRecord)
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
