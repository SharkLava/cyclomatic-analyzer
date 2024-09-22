package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FindCFiles recursively finds all .c and .h files in the given directory.
func FindCFiles(dir string) ([]string, error) {
	var cFiles []string
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(p, ".c") || strings.HasSuffix(p, ".h")) {
			cFiles = append(cFiles, p)
		}
		return nil
	})
	return cFiles, err
}

// ReadFile reads the content of a file.
func ReadFile(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WriteFile writes data to a file.
func WriteFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}
