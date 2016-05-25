package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/mattn/go-shellwords"
)

var dir string
var cmd string
var workers int
var recursive bool

const fileNameSymb = "{filename}"
const filePathSymb = "{filepath}"

func parseFlags() {
	_dir := flag.String("dir", "", "Dir containing files to run 'cmd' on.")
	_cmd := flag.String("cmd", "", fmt.Sprintf("Command to run, %s and %s will be substituded for files in 'dir'.", filePathSymb, fileNameSymb))
	_workers := flag.Int("workers", runtime.NumCPU(), "Number of go routines to spawn cmd on.")
	_recursive := flag.Bool("recursive", false, "Whether 'dir' should be searched recursively.")

	flag.Parse()

	workers = *_workers
	dir = *_dir
	recursive = *_recursive

	if len(*_cmd) == 0 {
		flag.Usage()
		os.Exit(0)
		return
	}
	cmd = *_cmd

}

func getFiles() []string {
	files := make([]string, 0, 512)

	var _getFiles func(string) []string
	_getFiles = func(dir string) []string {
		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Println(err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				_getFiles(path.Join(dir, entry.Name()))
			} else {
				files = append(files, path.Join(dir, entry.Name()))
			}
		}
		return files
	}

	return _getFiles(dir)
}

func getCmdArr(filePath string) []string {
	_cmd := strings.Replace(cmd, fileNameSymb, path.Base(filePath), -1)
	_cmd = strings.Replace(_cmd, filePathSymb, path.Join(dir, filePath), -1)
	cmdArr, err := shellwords.Parse(fmt.Sprintf("sh -c \"%s\"", _cmd))
	if err != nil {
		log.Fatal(err)
	}
	return cmdArr
}

func main() {
	parseFlags()

	fnChan := make(chan string, workers)

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for fileName := range fnChan {
				cmdArr := getCmdArr(fileName)
				cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("File name \"%s\":\n\t%s\n", fileName, out.String())
			}
			wg.Done()
		}()
	}

	go func() {
		for _, fileName := range getFiles() {
			fnChan <- fileName
		}
		close(fnChan)
	}()
	wg.Wait()
}
