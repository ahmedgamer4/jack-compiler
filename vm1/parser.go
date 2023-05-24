package main

import (
	"strings"
)

type CommandType int

const (
	C_ARITHMETIC CommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF_GOTO
	C_FUNCTION
	C_RETURN
	C_CALL
)

type CommandRecord struct {
	command string
	arg0    string
	arg1    string
}

func GetCommandType(command string) CommandType {
	typesRecord := map[string]CommandType{
		"add":      C_ARITHMETIC,
		"sub":      C_ARITHMETIC,
		"neg":      C_ARITHMETIC,
		"eq":       C_ARITHMETIC,
		"gt":       C_ARITHMETIC,
		"lt":       C_ARITHMETIC,
		"and":      C_ARITHMETIC,
		"or":       C_ARITHMETIC,
		"not":      C_ARITHMETIC,
		"push":     C_PUSH,
		"pop":      C_POP,
		"label":    C_LABEL,
		"goto":     C_GOTO,
		"if-goto":  C_IF_GOTO,
		"function": C_FUNCTION,
		"return":   C_RETURN,
		"call":     C_CALL,
	}
	return typesRecord[command]
}

func removeComments(line string) string {
	if slashIdx := strings.Index(line, "/"); slashIdx != -1 {
		return strings.TrimSpace(line[:slashIdx])
	}
	return strings.TrimSpace(line)
}

func ParseCommand(line string) CommandRecord {
	insArr := strings.Split(line, " ")
	res := CommandRecord{}
	res.command = strings.TrimSpace(insArr[0])

	switch GetCommandType(res.command) {
	case C_ARITHMETIC, C_RETURN:

	case C_PUSH, C_POP, C_CALL, C_FUNCTION:
		res.arg0 = strings.TrimSpace(insArr[1])
		res.arg1 = strings.TrimSpace(insArr[2])
	case C_LABEL, C_GOTO, C_IF_GOTO:
		res.arg1 = strings.TrimSpace(insArr[1])
	}

	return res
}
