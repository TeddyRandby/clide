package node

import (
	"clide/path"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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

type CommandNodeParam struct {
	Name string
	Type string
}

func (n CommandNode) Title() string       { return fmt.Sprintf("[%s] %s", n.Type, n.Name) }
func (n CommandNode) Description() string { return n.Path }
func (n CommandNode) FilterValue() string { return n.Path }

func (n CommandNode) Parameters() []CommandNodeParam {
	params := make([]CommandNodeParam, 0)
	steps := filepath.SplitList(n.Path)

    for _, step := range steps {
        if strings.Contains(step, "[") {
            params = append(params, CommandNodeParam{
                Name: strings.Trim(step, "[]"),
                Type: CommandNodeParamTypeInput,
            })
        } else if strings.Contains(step, "{") {
            params = append(params, CommandNodeParam{
                Name: strings.Trim(step, "{}"),
                Type: CommandNodeParamTypeSelect,
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

	if path.PathIsParam(pth) && path.PathIsLeaf(pth) {
		return nil, errors.New("A leaf cannot require a parameter")
	}

	node := new(CommandNode)
	node.Parent = parent
	node.Name = name
	node.Path = pth

	if path.PathIsLeaf(pth) {
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
	childs, err := path.PathChildren(dir)

	if err != nil {
		return nil, err
	}

	var nodes []CommandNode
	for _, child := range childs {
		if path.PathIsParam(child) {
			grandchilds, err := children(parent, child)

			if err != nil {
				return nil, err
			}

			nodes = append(nodes, grandchilds...)
		} else {
			node, err := New(parent, child)

			if err != nil {
				return nil, err
			}

			if node != nil {
				nodes = append(nodes, *node)
			}
		}
	}

	return nodes, nil
}
