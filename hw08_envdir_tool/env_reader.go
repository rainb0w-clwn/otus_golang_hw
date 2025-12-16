package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	m := make(Environment)
	w := sync.WaitGroup{}
	mu := &sync.Mutex{}
	for _, entry := range files {
		w.Add(1)
		go func(_entry os.DirEntry) {
			defer w.Done()
			if !_entry.IsDir() && !strings.ContainsRune(_entry.Name(), '=') {
				file, err := os.Open(filepath.Join(dir, _entry.Name()))
				if err == nil {
					r := bufio.NewReader(file)
					lineB0, err := r.ReadBytes('\n')
					lineB1 := bytes.TrimRight(lineB0, " \t\n")
					lineB2 := bytes.ReplaceAll(lineB1, []byte{0x00}, []byte{'\n'})
					mu.Lock()
					m[_entry.Name()] = EnvValue{
						Value:      string(lineB2),
						NeedRemove: len(lineB0) == 0 && err != nil,
					}
					mu.Unlock()
				}
			}
		}(entry)
	}
	w.Wait()
	return m, nil
}
