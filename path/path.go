package path

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	ParamInputPrefix  = "["
	ParamInputSuffix  = "]"
	ParamSelectPrefix = "{"
	ParamSelectSuffix = "}"
	ParamBracketChars = ParamInputPrefix + ParamInputSuffix + ParamSelectPrefix + ParamSelectSuffix
	ParamChars        = ParamBracketChars
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func HasSibling(path string, child string) string {
	sibling := filepath.Join(path, "..", child)

	if Exists(sibling) {
		return sibling
	}

	return ""
}

func Children(path string) ([]string, error) {
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

func Parent(path string) string {
	return filepath.Join(path, "..")
}

func IsLeaf(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !info.IsDir()
}

func hasPrefixAndSuffix(s string, prefix string, suffix string) bool {
	return strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix)
}

func IsSelectParameter(path string) bool {
	p := filepath.Base(path)

	return hasPrefixAndSuffix(p, ParamSelectPrefix, ParamSelectSuffix)
}

func IsInputParameter(path string) bool {
	p := filepath.Base(path)

	return hasPrefixAndSuffix(p, ParamInputPrefix, ParamInputSuffix)
}

func IsParameter(path string) bool {
	return IsSelectParameter(path) || IsInputParameter(path)
}

func IsModule(path string) bool {
	return !IsLeaf(path)
}

func IsRoot(path string) bool {
	return Exists(filepath.Join(path, ".clide"))
}

func findRoot(path string) (string, error) {
	if Exists(filepath.Join(path, ".git")) {
		if Exists(filepath.Join(path, ".clide")) {
			return filepath.Join(path, ".clide"), nil
		}

		return "", nil
	}

	parentPath := filepath.Join(path, "..")

	if parentPath == path {
		return "", nil
	}

	return findRoot(parentPath)
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
