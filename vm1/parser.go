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
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type CommandRecord struct {
	command string
	segment string
	num     string
}

func GetCommandType(command string) CommandType {
	typesRecord := map[string]CommandType{
		"add":  C_ARITHMETIC,
		"sub":  C_ARITHMETIC,
		"neg":  C_ARITHMETIC,
		"eq":   C_ARITHMETIC,
		"gt":   C_ARITHMETIC,
		"lt":   C_ARITHMETIC,
		"and":  C_ARITHMETIC,
		"or":   C_ARITHMETIC,
		"not":  C_ARITHMETIC,
		"push": C_PUSH,
		"pop":  C_POP,
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
	if GetCommandType(insArr[0]) == C_ARITHMETIC {
		res.command = removeComments(insArr[0])
	}
	if GetCommandType(insArr[0]) == C_PUSH || GetCommandType(insArr[0]) == C_POP {
		res.command = removeComments(insArr[0])
		res.segment = removeComments(insArr[1])
		res.num = removeComments(insArr[2])
	}
	return res
}
