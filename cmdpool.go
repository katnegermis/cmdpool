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

const fileNameSymb = "{filename}"
const filePathSymb = "{filepath}"

func parseFlags() {
	_dir := flag.String("dir", "", "Path to files to run command on.")
	_cmd := flag.String("cmd", "", fmt.Sprintf("Command to run, %s and %s will be substituded.", filePathSymb, fileNameSymb))
	_workers := flag.Int("workers", runtime.NumCPU(), "Number of go routines to spawn cmd on.")

	flag.Parse()

	workers = *_workers

	if len(*_cmd) == 0 {
		flag.Usage()
		os.Exit(0)
		return
	}
	cmd = *_cmd

	dir = *_dir
}

func getFiles() []string {
	files := make([]string, 0, 512)
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files
}

func getCmdArr(fileName string) []string {
	_cmd := strings.Replace(cmd, fileNameSymb, fileName, -1)
	_cmd = strings.Replace(_cmd, filePathSymb, path.Join(dir, fileName), -1)
	cmdArr, err := shellwords.Parse(_cmd)
	if err != nil {
		log.Fatal(err)
	}
	return cmdArr
}

func main() {
	parseFlags()

	fnChan := make(chan string)

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