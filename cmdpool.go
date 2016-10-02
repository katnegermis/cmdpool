package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	shellwords "github.com/mattn/go-shellwords"
)

var (
	cmd       string
	workers   int
	recursive bool
	num       int
	files     []string
	print     bool
)

const fileNameSymb = "{filename}"
const filePathSymb = "{filepath}"
const iterationSymb = "{iteration}"

func parseFlags() {
	flag.IntVar(&num, "num", 1, "Number of times to run 'cmd'")
	flag.StringVar(&cmd, "cmd", "", fmt.Sprintf("Command to run, %s and %s will be substituded for files in 'dir'.", filePathSymb, fileNameSymb))
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of go routines to spawn cmd on.")
	flag.BoolVar(&print, "print", true, "Whether to print results of 'cmd' to stdout.")

	flag.Parse()

	files = flag.Args()
	files = files

	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(0)
		return
	}

}

type workDef struct {
	num int
	cmd string
}

func makeShellCommand(w workDef) []string {
	_cmd := strings.Replace(cmd, fileNameSymb, path.Base(w.cmd), -1)
	_cmd = strings.Replace(_cmd, filePathSymb, w.cmd, -1)
	_cmd = strings.Replace(_cmd, iterationSymb, strconv.Itoa(w.num), -1)
	cmdArr, err := shellwords.Parse(fmt.Sprintf("sh -c \"%s\"", _cmd))
	if err != nil {
		log.Fatal(err)
	}
	return cmdArr
}

func main() {
	parseFlags()

	workChan := make(chan workDef, workers*2)

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			for w := range workChan {
				args := makeShellCommand(w)
				cmd := exec.Command(args[0], args[1:]...)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					fmt.Println("error: ", err)
				}
				if print {
					fmt.Printf("`%s`: %s\n", strings.Join(args, " "), out.String())
				}
			}
			wg.Done()
		}()
	}

	if len(files) > 0 {
		for i := 0; i < num; i++ {
			for _, f := range files {
				workChan <- workDef{num: i, cmd: f}
			}
		}
	} else {
		for i := 0; i < num; i++ {
			workChan <- workDef{num: i, cmd: cmd}
		}
	}
	close(workChan)

	wg.Wait()
}
