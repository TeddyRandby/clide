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

func pathChildren(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var children []string

	for _, child := range files {
		children = append(children, child.Name())
	}

	return children, nil
}

func PathIsLeaf(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !info.IsDir()
}

func PathIsParam(path string) bool {
    return strings.ContainsAny(path, "(){}[]")
}

func PathIsDir(path string) bool {
    return !PathIsLeaf(path) 
}

const (
	NodeTypeCommand  = "Command"
	NodeTypeModule   = "Module"
	NodeTypeParameter = "Parameter"
)

type CommandNode struct {
	Name string
	Type string
    children []CommandNode
}

func (n CommandNode) Title() string       { return n.Name }
func (n CommandNode) Description() string { return n.Type }
func (n CommandNode) FilterValue() string { return n.Name }

func PathChoices(path string) ([]CommandNode, error) {
	children, err := pathChildren(path)

	if err != nil {
		return nil, err
	}

	var nodes []CommandNode
	for _, child := range children {
		if !strings.HasPrefix(child, ".") {
			if PathIsLeaf(filepath.Join(path, child)) {
				nodes = append(nodes, CommandNode{Name: child, Type: NodeTypeCommand})
			} else if PathIsParam(filepath.Join(path, child)) {
				nodes = append(nodes, CommandNode{Name: child, Type: NodeTypeParameter})
			} else {
				nodes = append(nodes, CommandNode{Name: child, Type: NodeTypeModule})
			}
		}
	}

	return nodes, nil
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
