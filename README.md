Run commands in your shell over a directory of files, using a pool of workers.

Usage
=====
~~~
Usage of cmdpool:
  -cmd string
        Command to run, {filepath}, {filename}, {iteration} will be substituded for 'files' and a unique number for each execution
  -num int
        Number of times to run 'cmd' (default 1)
  -print
        Whether to print results of 'cmd' to stdout. (default true)
  -workers int
        Number of go routines to spawn cmd on. (default 4)
  [files] string
        A list of file names to be given as input to 'cmd'
~~~

Examples
========

~~~
./cmdpool -workers 3 -cmd "echo {iteration}: {filename}; sleep 3" README.md cmdpool.go
~~~
