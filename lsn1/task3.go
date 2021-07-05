package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CreateFiles(dir string, countFiles int) error {
	if err := os.Mkdir(dir, 0700); err != nil { //os.ModeDir
		fmt.Println("Failed to make directory :", err)
		return err
	}
	for i := 0; i < countFiles; i++ {
		fileName := fmt.Sprintf("%s/file_%d", dir, i) //dir + "/file_" + strconv.Itoa(i)
		fo, _ := os.Create(fileName)
		defer fo.Close()
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %s: %s\n", dir, err)
	}
	fmt.Printf("Created %d files.\n", len(files))
	fmt.Println("Cleaning up")
	os.RemoveAll(dir)
	fmt.Println("Done")
	return nil
}

// todo: use https://yourbasic.org/golang/temporary-file-directory/
