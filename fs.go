package main

import (
	"os"
	"strings"
)

func listFiles(dir string) []string {
	fs, err := os.ReadDir(dir)
	check(err)
	var files []string
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue // skip hidden files, or files currently being copied
		}
		files = append(files, f.Name())
	}
	return files
}
