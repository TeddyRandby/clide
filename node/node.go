package node

import (
	"clide/path"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

const (
	NodeTypeCommand = "Command"
	NodeTypeModule  = "Module"
)

type CommandNode struct {
	Name     string
	Path     string
	Type     string
	Children []CommandNode
	Parent   *CommandNode
}

const (
	CommandNodeParamTypeInput    = "Input"
	CommandNodeParamTypeSelect   = "Select"
	CommandNodeParamTypeCheckbox = "Checkbox"
)

type CommandNodeParameters struct {
	Shortcut string
	Name     string
	Type     string
}

func (n CommandNode) Title() string       { return fmt.Sprintf("[%s] %s", n.Type, n.Name) }
func (n CommandNode) Description() string { return n.Path }
func (n CommandNode) FilterValue() string { return n.Path }

func nameAndShortcut(original string) (string, string) {
	name := strings.ToLower(strings.Trim(original, "[]{}<>"))

	var shortcut string
	for _, char := range original {
		if unicode.IsUpper(char) {
			shortcut += string(char)
		}
	}

	return name, shortcut
}

func (n CommandNode) Parameters() []CommandNodeParameters {
	params := make([]CommandNodeParameters, 0)

	steps := strings.Split(n.Path, "/")

	root := slices.Index(steps, ".clide")

	steps = steps[root+1:]

	for _, step := range steps {
		if strings.Contains(step, "[") {
			name, shortcut := nameAndShortcut(step)
			params = append(params, CommandNodeParameters{
				Name:     name,
				Shortcut: shortcut,
				Type:     CommandNodeParamTypeInput,
			})
		} else if strings.Contains(step, "{") {
			name, shortcut := nameAndShortcut(step)
			params = append(params, CommandNodeParameters{
				Name:     name,
				Shortcut: shortcut,
				Type:     CommandNodeParamTypeSelect,
			})
		} else if strings.Contains(step, "<") {
			name, shortcut := nameAndShortcut(step)
			params = append(params, CommandNodeParameters{
				Name:     name,
				Shortcut: shortcut,
				Type:     CommandNodeParamTypeCheckbox,
			})
		}
	}

	return params
}

func Root() (*CommandNode, error) {
	root, err := path.FindRoot()

	if err != nil {
		return nil, err
	}

	if root == "" {
		return nil, err
	}

	node, err := New(nil, root)

	if err != nil {
		return nil, err
	}

	return node, nil
}

func New(parent *CommandNode, pth string) (*CommandNode, error) {
	name := filepath.Base(pth)

	if strings.HasPrefix(name, ".") && name != ".clide" {
		return nil, nil
	}

	if path.IsParameter(pth) && path.IsLeaf(pth) {
		return nil, errors.New("A leaf cannot require a parameter")
	}

	node := new(CommandNode)
	node.Parent = parent
	node.Name = name
	node.Path = pth

	if path.IsLeaf(pth) {
		node.Type = NodeTypeCommand
		return node, nil
	}

	node.Type = NodeTypeModule

	childs, err := children(node, pth)

	if err != nil {
		return nil, err
	}

	node.Children = childs
	return node, nil
}

func children(parent *CommandNode, dir string) ([]CommandNode, error) {
	childs, err := path.Children(dir)

	if err != nil {
		return nil, err
	}

	var nodes []CommandNode
	for _, child := range childs {
		if path.IsParameter(child) {
			grandchilds, err := children(parent, child)

			if err != nil {
				return nil, err
			}

			nodes = append(nodes, grandchilds...)
		} else {
			if path.IsLeaf(child) {
				if strings.Contains(filepath.Base(child), ".") {
					node, err := New(parent, child)

					if err != nil {
						return nil, err
					}

					if node != nil {
						nodes = append(nodes, *node)
					}
				}
			}

			if path.IsModule(child) {
				node, err := New(parent, child)

				if err != nil {
					return nil, err
				}

				if node != nil {
					nodes = append(nodes, *node)
				}
			}

		}
	}

	return nodes, nil
}
