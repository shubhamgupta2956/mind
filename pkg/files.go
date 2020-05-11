package pkg

import (
	"fmt"
	"os"
)

// CreateDirIfNotExist creates a directory if it doesn't exist else does nothing.
func CreateDirIfNotExist(dir string) error {
	exists, info, err := isExist(dir)
	if err != nil {
		return err
	}

	if exists {
		if !info.IsDir() {
			return fmt.Errorf("path '%s' exists and is not a directory", dir)
		}

		return nil
	}

	return os.MkdirAll(dir, os.ModePerm)
}

// CreateFileIfNotExist creates a file if it doesn't exist else does nothing.
// Writes content to the empty file while creating.
func CreateFileIfNotExist(file, content string) error {
	exists, info, err := isExist(file)
	if err != nil {
		return err
	}

	if exists {
		if info.IsDir() {
			return fmt.Errorf("path '%s' exists and is not a file", file)
		}

		return nil
	}

	newFile, err := os.Create(file)
	if err != nil {
		return err
	}

	if _, err := newFile.WriteString(content); err != nil {
		return err
	}

	return newFile.Close()
}

// isExist tells if the file exists or not (whether it is dir or file).
func isExist(file string) (bool, os.FileInfo, error) {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	return true, info, nil
}
