package model

import (
	"clide/node"
	"clide/path"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

func (m Clide) Backtrack() Clide {
	parent := m.node.Parent

	if parent == nil {
		return m
	}

	return m.PromptPath(parent)
}

func (m Clide) Root() Clide {
	return m.PromptPath(m.root)
}

func (m Clide) Error(err string) Clide {
	return Clide{
		width:  m.width,
		height: m.height,
		root:   m.root,
		params: m.params,
		param:  m.param,
		args:   m.args,
		keymap: m.keymap,
		help:   m.help,
		ready:  true,
		state:  ClideStateError,
		error:  err,
	}
}

func (m Clide) Command(n *node.CommandNode) Clide {
	m.node = n

	if n.Type != node.NodeTypeCommand {
		return m.Error("Can only execute commands")
	}

	m.params = n.Parameters()

	if len(m.params) > m.param {
		return m.nextArgument()
	}

	return m.execute()
}

func (m Clide) execute() Clide {
	cmd := exec.Command(m.node.Path)

	output, err := cmd.Output()

	if err != nil {
		return m.Error(err.Error())
	}

	return m.Done(string(output))
}

func (m Clide) Done(message string) Clide {
	clide := Clide{
		ready:    m.ready,
		width:    m.width,
		height:   m.height,
		node:     m.node,
		root:     m.root,
		args:     m.args,
		keymap:   m.keymap,
		help:     m.help,
		params:   nil,
		param:    0,
		state:    ClideStateDone,
		viewport: viewport.New(m.width, m.height),
	}

	clide.viewport.SetContent(message)
	return clide
}

func (m Clide) SetAndPromptNextArgument(value string) Clide {
	os.Setenv(m.params[m.param].Name, value)
	m.param++

	if len(m.params) > m.param {
		return m.nextArgument()
	}

	return m.execute()
}

func (m Clide) nextArgument() Clide {
	param := m.params[m.param]

	shortcutValue := m.args[param.Shortcut]
	if shortcutValue != "" {
		return m.SetAndPromptNextArgument(shortcutValue)
	}

	switch param.Type {
	case node.CommandNodeParamTypeInput:
		return m.PromptInput()
	case node.CommandNodeParamTypeSelect:
		return m.PromptSelect()
	}

	return m.Error("Invalid parameter type")
}

func (m Clide) PromptPath(n *node.CommandNode) Clide {
	m.node = n

	if m.node == nil {
		return m.Error("Invalid node")
	}

	options := m.node.Children

	if options == nil || len(options) == 0 {
		return m.Error(fmt.Sprintf("No commands found in %s", m.node.Path))
	}

	items := make([]list.Item, len(options))

	for i, choice := range options {
		items[i] = list.Item(choice)
	}

	c := Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		node:   m.node,
		root:   m.root,
		params: m.params,
		param:  m.param,
		args:   m.args,
		keymap: m.keymap,
		help:   m.help,
		state:  ClideStatePathSelect,
		list:   list.New(items, list.NewDefaultDelegate(), m.width, m.height),
	}

	return c
}

type item struct {
	name, desc string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.name }

func (m Clide) PromptSelect() Clide {
	name := m.params[m.param].Name

	sibling := path.HasSibling(m.node.Path, name)

	if sibling == "" {
		return m.Error(fmt.Sprintf("Invalid parameter: No %s found in %s", name, path.Parent(m.node.Path)))
	}

	cmd := exec.Command(sibling)

	output, err := cmd.Output()

	if err != nil {
		return m.Error(err.Error())
	}

	options := strings.Split(string(output), "\n")

	if len(options) == 0 {
		return m.Error(fmt.Sprintf("%s yielded no options", name))
	}

	items := make([]list.Item, len(options))

	for i, choice := range options {
		values := strings.Split(choice, ":")
		items[i] = list.Item(item{values[0], values[1]})
	}

	return Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		node:   m.node,
		root:   m.root,
		params: m.params,
		param:  m.param,
		args:   m.args,
		keymap: m.keymap,
		help:   m.help,
		state:  ClideStatePromptSelect,
		list:   list.New(items, list.NewDefaultDelegate(), m.width, m.height),
	}
}

func (m Clide) PromptInput() Clide {
	c := Clide{
		ready:     m.ready,
		width:     m.width,
		height:    m.height,
		node:      m.node,
		root:      m.root,
		params:    m.params,
		param:     m.param,
		args:      m.args,
		keymap:    m.keymap,
		help:      m.help,
		state:     ClideStatePromptInput,
		textinput: textinput.New(),
	}

	c.textinput.Width = m.width
	c.textinput.CharLimit = 50
	c.textinput.Placeholder = "Enter a value"

	c.textinput.Focus()

	return c
}
