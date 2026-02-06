package library

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func Scan(path string) ([]string, error) {

	var files []string

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("Failed to access %v: %v", path, err)
		}

		if d.IsDir() {
			return nil
		}

		if strings.EqualFold(filepath.Ext(d.Name()), ".mp3") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Something went wrong %v", err)
	}

	return files, nil
}
