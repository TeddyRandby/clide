package path

import (
	"os"
	"path/filepath"
	"strings"
)

func PathExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func PathChildren(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var children []string

	for _, child := range files {
		children = append(children, filepath.Join(path, child.Name()))
	}

	return children, nil
}

func PathParent(path string) string {
    return filepath.Join(path, "..")
}

func PathIsLeaf(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !info.IsDir()
}

func PathIsParam(path string) bool {
    return strings.ContainsAny(filepath.Base(path), "(){}[]")
}

func PathIsDir(path string) bool {
    return !PathIsLeaf(path) 
}

func PathIsRoot(path string) bool {
    return PathExists(filepath.Join(path, ".clide"))
}


func findRoot(path string) (string, error) {
	if PathExists(filepath.Join(path, ".git")) {
		if PathExists(filepath.Join(path, ".clide")) {
			return filepath.Join(path, ".clide"), nil
		}

		return "", nil
	}

	return findRoot(filepath.Join(path, ".."))
}

func FindRoot() (string, error) {
	path, err := os.Getwd()

	if err != nil {
		return "", err
	}

	root, err := findRoot(path)

	if err != nil {
		return "", err
	}

	if root == "" {
		return "", nil
	}

	return root, nil
}
