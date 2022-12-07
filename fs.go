package main

import "os"

func listFiles(dir string) []string {
	fs, err := os.ReadDir(dir)
	check(err)
	var files []string
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		files = append(files, f.Name())
	}
	return files
}
