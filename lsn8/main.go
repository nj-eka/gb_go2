package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"hash/adler32"
	"io"
	"os"
	"os/user"
	"runtime"
	"sync"
)

type FilePathHash struct {
	FilePath string
	Hash     uint32
}

func IterFiles(startPath string, filePathChan chan string) {
	entries, err := os.ReadDir(startPath)
	if err != nil {
		log.Panic(err)
	}

	for _, fileName := range entries {
		curPath := fmt.Sprintf("%s/%s", startPath, fileName.Name())
		if fileName.IsDir() {
			IterFiles(curPath, filePathChan)
		} else {
			filePathChan <- curPath
		}

	}
}

var (
	rootDir = flag.String("root", ".", "root directory")
	removeDups = flag.Bool("rm", false, "remove duplicates")
)

func initLog() error{

}

func main() {
	flag.Parse()
	initLog()
	log.SetFormatter(&log.TextFormater{})
	logFields := log.Fields{
		"pid": os.Getpid(),
		"user": user.Current().Username,
	}
	logger := log.WithFields(logFields)

	standardFields := log.Fields{
		"host": "srv42",
	}

	var wg sync.WaitGroup
	filePathChan := make(chan string)
	fileHashChan := make(chan *FilePathHash)

	go func() {
		defer close(filePathChan)
		IterFiles(*rootDir, filePathChan)
	}()

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for filePath := range filePathChan {
				file, err := os.Open(filePath)
				if err != nil {
					log.Println(err)
					continue
				}

				hash := adler32.New()
				io.Copy(hash, file)
				file.Close()

				fileHashChan <- &FilePathHash{
					FilePath: filePath,
					Hash:     hash.Sum32(),
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(fileHashChan)
	}()

	copies := make(map[uint32][]string)

	for fileHash := range fileHashChan {
		filesPath := copies[fileHash.Hash]
		filesPath = append(filesPath, fileHash.FilePath)
		copies[fileHash.Hash] = filesPath
	}

	if *removeDups {
		for _, paths := range copies {
			for i := 1; i < len(paths); i++ {
				err := os.Remove(paths[i])
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println("file: ", paths[i], " - removed")
				}
			}
		}
	}

}
