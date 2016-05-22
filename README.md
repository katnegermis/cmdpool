Run commands in your shell using a pool of workers.

Usage
=====
~~~
Usage of ./cmdpool:
  -cmd string
        Command to run, {filepath} and {filename} will be substituded.
  -dir string
        Path to files to run command on.
  -workers int
        Number of go routines to spawn cmd on. (default 4)
~~~

Examples
~~~
./cmdpool -workers 50 -path /tmp -cmd "./tool {filepath}"
~~~
