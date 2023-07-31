package node

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/TeddyRandby/clide/path"
	"golang.org/x/exp/slices"
)

const (
	NodeTypeCommand = " "
	NodeTypeModule  = " "
)

type CommandNode struct {
	Name     string
	Shortcut string
	Path     string
	Type     string
	Children []CommandNode
	Parent   *CommandNode
}

const (
	CommandNodeParamTypeInput  = "Input"
	CommandNodeParamTypeSelect = "Select"
)

type CommandNodeParameters struct {
	Shortcut string
	Name     string
	Type     string
	Value    string
}

func (n CommandNode) Title() string {
	if n.Shortcut != "" {
		return fmt.Sprintf("%s %s (%s)", n.Type, n.Name, n.Shortcut)
	}
	return fmt.Sprintf("%s %s", n.Type, n.Name)
}

func (n CommandNode) Description() string { return n.clideRelativePath() }

func (n CommandNode) FilterValue() string { return n.Name }

func moduleNameAndShortcut(original string) (string, string) {
	name := strings.Split(original, ".")[0]

	if name == "" {
		name = original
	}

	var shortcut string
	for _, char := range original {
		if unicode.IsUpper(char) {
			shortcut += string(char)
		}
	}

	return strings.ToLower(name), strings.ToLower(shortcut)
}

func parameterNameAndShortcut(original string) (string, string) {
	name := strings.Trim(original, path.ParamChars)

	var shortcut string
	for _, char := range original {
		if unicode.IsUpper(char) {
			shortcut += string(char)
		}
	}

	return strings.ToLower(name), strings.ToLower(shortcut)
}

func (n CommandNode) Leaves() []CommandNode {
	leaves := make([]CommandNode, 0)

	for _, child := range n.Children {
		switch child.Type {
		case NodeTypeCommand:
			leaves = append(leaves, child)
		case NodeTypeModule:
			leaves = append(leaves, child.Leaves()...)
		}
	}

	return leaves
}

func (n CommandNode) Steps() string {
	steps := make([]string, 0)

	node := &n

	for node.Parent != nil {
		steps = append(steps, node.Name)
		node = node.Parent
	}

	result := make([]string, len(steps))
	// Reverse the steps slice
	for i := 0; i < len(steps); i++ {
		result[len(steps)-1-i] = steps[i]
	}

	return strings.Join(result, " ")
}

func (n CommandNode) clideRelativePath() string {
	return filepath.Join(n.clideRelativeSteps()...)
}

func (n CommandNode) clideRelativeSteps() []string {
	steps := strings.Split(n.Path, "/")

	root := slices.Index(steps, ".clide")

	return steps[root+1:]
}

func (n CommandNode) Parameters() []CommandNodeParameters {
	params := make([]CommandNodeParameters, 0)

	steps := n.clideRelativeSteps()

	for _, step := range steps {
		if strings.ContainsAny(step, path.ParamInputChars) {
			name, shortcut := parameterNameAndShortcut(step)
			params = append(params, CommandNodeParameters{
				Name:     name,
				Shortcut: shortcut,
				Type:     CommandNodeParamTypeInput,
			})
		} else if strings.ContainsAny(step, path.ParamSelectChars) {
			name, shortcut := parameterNameAndShortcut(step)
			params = append(params, CommandNodeParameters{
				Name:     name,
				Shortcut: shortcut,
				Type:     CommandNodeParamTypeSelect,
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
	original_name := filepath.Base(pth)

	if strings.HasPrefix(original_name, ".") && original_name != ".clide" {
		return nil, nil
	}

	if path.IsParameter(pth) && path.IsLeaf(pth) {
		return nil, errors.New("A leaf cannot require a parameter")
	}

	name, shortcut := moduleNameAndShortcut(original_name)

	node := new(CommandNode)
	node.Parent = parent
	node.Name = name
	node.Shortcut = shortcut
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

func (n CommandNode) findChild(name string) (*CommandNode, error) {
	for _, child := range n.Children {
		if child.Name == name {
			return &child, nil
		}
	}
	return nil, errors.New("Child not found")
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

	for i, a := range nodes {
		for j, b := range nodes {
			if i != j {
				if a.Name == b.Name {
					return nil, errors.New(fmt.Sprintf("Duplicate leaf %s", a.Name))
				}
			}
		}
	}

	return nodes, nil
}
