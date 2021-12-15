package file

import (
	"os"
	"path/filepath"
)

func Create(path string) (*os.File, error) {
	var f *os.File
	var err error

	if path == "" {
		f = os.Stdout
		return f, nil
	}

	dir, _ := filepath.Split(path)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(path)
}
