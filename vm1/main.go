package main

import (
	"fmt"
	"os"
)

func main() {
	lines, outputFileName := HandleOpen()

	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer outputFile.Close()

	outputFile.WriteString(Bootstrap())

	for _, line := range lines {
		lineWithoutComments := ParseCommand(line)
		ins := Translate(lineWithoutComments)
		outputFile.WriteString(ins)
	}
}
