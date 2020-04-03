package _file

import (
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

// check filename exists
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ExistsInPath(root string, fileName string) bool {
	fullPath, err := FullPath(fileName)
	if err != nil {
		return false
	}
	if !Exists(fullPath) {
		return false
	}
	fullRoot, err := FullPath(root)
	if err != nil {
		return false
	}
	if !Exists(fullRoot) {
		return false
	}
	return strings.HasPrefix(fullPath, fullRoot)
}

func FullPath(fileName string) (string, error) {
	path, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	return strings.Replace(path, "\\", "/", -1), nil
}
