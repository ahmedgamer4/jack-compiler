package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func HandleOpen() ([]string, string) {
	isDir := flag.Bool("d", false, "Compile a directory")
	flag.Parse()
	lines := make([]string, 0)
	outputFile := ""
	fmt.Println(*isDir)

	if *isDir {
		lines = openDir()
		outputFile = GetInArg() + "/" + path.Base(GetInArg()) + ".asm"
		fmt.Println(outputFile)
	} else {
		lines = openFile(GetInArg())
		outputFile = path.Dir(GetInArg()) + "/" + strings.Split(path.Base(GetInArg()), ".")[0] + ".asm"
	}
	return lines, outputFile
}

func openFile(file string) []string {
	readFile, err := os.Open(file)

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

func openDir() []string {
	files, err := ioutil.ReadDir(GetInArg())

	res := make([]string, 0)

	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".vm") {
			fileContent := openFile(GetInArg() + "/" + file.Name())
			res = append(res, fileContent...)
			fmt.Println(GetInArg() + "/" + file.Name())
		}
	}

	return res
}

func GetInArg() string {
	return os.Args[len(os.Args)-1]
}
