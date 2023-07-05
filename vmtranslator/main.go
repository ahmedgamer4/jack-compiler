package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	lines, outputFileName := HandleOpen()

	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer outputFile.Close()

	if !strings.HasSuffix(GetInArg(), ".vm") {
		outputFile.WriteString(Bootstrap())
	}

	for _, line := range lines {
		lineWithoutComments := ParseCommand(line)
		ins := Translate(lineWithoutComments)
		outputFile.WriteString(ins)
	}
}
