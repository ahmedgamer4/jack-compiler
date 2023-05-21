package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func OpenFile() []string {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as an arugment")
	}

	filePath := os.Args[1]
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	readFile.Close()

	return fileLines
}

func GetPathName() []string {
	return strings.Split(os.Args[1], "/")
}

func GetFilename() string {
	filePathSlice := GetPathName()

	inputFilename := filePathSlice[len(filePathSlice)-1]
	if inputFilename[len(inputFilename)-3:] != ".vm" {
		color.Red("File should be of format .vm")
		return ""
	}

	filename := strings.Split(inputFilename, ".")[0] + ".asm"
	return filename
}
