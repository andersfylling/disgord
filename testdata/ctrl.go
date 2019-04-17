package testdata

import (
	"os"
	"path/filepath"
	"strings"
)

func GetDataForDir(dir string, filter string) (files []string, err error) {
	root := "testdata/" + dir
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path[len(root):], filter) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := range files {
		if files[i] == root {
			files[i] = files[len(files)-1]
			files = files[:len(files)-1]
			break
		}
	}

	return files, nil
}
