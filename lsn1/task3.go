package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func createFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", fileName)
	}
	defer file.Close()
}

func CreateFiles(dir string, countFiles int) error {
	if err := os.Mkdir(dir, os.ModeDir|os.ModePerm); err != nil {
		fmt.Println("Failed to make directory :", err)
		return err
	}
	for i := 0; i < countFiles; i++ {
		fileName := fmt.Sprintf("%s/file_%d", dir, i)
		createFile(fileName)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %s: %s\n", dir, err) // too many open files
	}
	fmt.Printf("Created %d files.\n", len(files))
	fmt.Println("Cleaning up")
	os.RemoveAll(dir)
	fmt.Println("Done")
	return nil
}

// todo: use https://yourbasic.org/golang/temporary-file-directory/
