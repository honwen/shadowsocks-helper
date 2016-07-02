package main

import (
	"os"
	"path/filepath"
	"regexp"
)

func getCurrDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func getFileList(path string, checkFunc func(string) bool) (files []string) {
	filepath.Walk(path, func(path string, fInfo os.FileInfo, err error) error {
		if fInfo == nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		if checkFunc == nil || checkFunc(path) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func isGuiConfig(s string) bool {
	return regexp.MustCompile(`gui-config.*\.json$`).MatchString(s)
}
