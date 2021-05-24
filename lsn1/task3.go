package main

import (
	"os"
	"strconv"
)

func createFiles(dir string, countFiles int) {
	os.Mkdir(dir, os.ModeDir)
	for i := 0; i < countFiles; i++ {
		fileName := dir + "/file" + strconv.Itoa(i)
		fo, _ := os.Create(fileName)
		defer fo.Close()
	}
}

func main() {
	createFiles(".tmp", 10000)
}

// todo: use https://yourbasic.org/golang/temporary-file-directory/
