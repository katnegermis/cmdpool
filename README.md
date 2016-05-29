Run commands in your shell over a directory of files, using a pool of workers.

Usage
=====
~~~
Usage of ./cmdpool:
  -cmd string
        Command to run, {filepath} and {filename} will be substituded for files in 'dir'.
  -dir string
        Dir containing files to run 'cmd' on.
  -recursive
        Whether 'dir' should be searched recursively.
  -workers int
        Number of go routines to spawn cmd on. (defaults to number of logical cores)
~~~

Examples
========

~~~
./cmdpool -workers 3 -dir /tmp -cmd "echo {filename}; sleep 3"
~~~
