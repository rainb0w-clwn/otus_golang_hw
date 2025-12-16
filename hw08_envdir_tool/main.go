package main

import (
	"os"
)

const ExitCode = 111

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		os.Exit(ExitCode)
	}
	env, err := ReadDir(args[0])
	if err != nil {
		os.Exit(ExitCode)
	}
	r := RunCmd(args[1:], env)
	os.Exit(r)
}
