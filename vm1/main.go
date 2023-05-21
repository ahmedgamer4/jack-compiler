package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	pathWithAsm := strings.Join(GetPathName()[0:len(GetPathName())-1], "/") + "/" + GetFilename()
	outputFile, err := os.OpenFile(pathWithAsm, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer outputFile.Close()

	lines := OpenFile()

	for _, line := range lines {
		lineWithoutComments := ParseCommand(line)
		ins := Translate(lineWithoutComments)
		if ins != "/" {
			outputFile.WriteString(ins)
		}
	}
}
