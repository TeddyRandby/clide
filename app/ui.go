package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TeddyRandby/clide/node"
	"github.com/TeddyRandby/clide/path"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	delegate = func() list.DefaultDelegate {
		d := list.NewDefaultDelegate()
		d.Styles.DimmedTitle.Foreground(white)
		d.Styles.DimmedDesc.Foreground(gray)
		d.Styles.NormalTitle.Foreground(fg)
		d.Styles.NormalDesc.Foreground(white)
		d.Styles.SelectedTitle.Foreground(purple)
		d.Styles.SelectedDesc.Foreground(white)
		d.Styles.SelectedTitle.BorderForeground(purple)
		d.Styles.SelectedDesc.BorderForeground(purple)
		return d
	}()
)

func (m Clide) Backtrack() (Clide, tea.Cmd) {
	parent := m.node.Parent

	if parent == nil {
		return m, nil
	}

	return m.PromptPath(parent)
}

func (m Clide) Root() (Clide, tea.Cmd) {
	return m.PromptPath(m.root)
}

func (m Clide) Error(err string) (Clide, tea.Cmd) {
	return Clide{
		width:  m.width,
		height: m.height,
		root:   m.root,
		node:   m.node,
		params: m.params,
		param:  m.param,
		args:   m.args,
		keymap: m.keymap,
		help:   m.help,
		ready:  m.ready,
		state:  ClideStateError,
		error:  err,
	}, nil
}

func (m Clide) Command(n *node.CommandNode) (Clide, tea.Cmd) {
	m.node = n

	if n.Type != node.NodeTypeCommand {
		return m.Error("Can only execute commands")
	}

	m.params = n.Parameters()

	if len(m.params) > m.param {
		return m.nextParameter()
	}

	return m.Done()
}

func (m Clide) Done() (Clide, tea.Cmd) {
	clide := Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		node:   m.node,
		root:   m.root,
		args:   m.args,
		keymap: m.keymap,
		help:   m.help,
		params: nil,
		param:  0,
		state:  ClideStateDone,
	}

	return clide, tea.Quit
}

func (m Clide) SetAndPromptNextArgument(value string) (Clide, tea.Cmd) {
	os.Setenv(m.params[m.param].Name, value)
	m.param++

	if len(m.params) > m.param {
		return m.nextParameter()
	}

	return m.Done()
}

func (m Clide) nextParameter() (Clide, tea.Cmd) {
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

func (m Clide) newlist(items []list.Item) list.Model {

	l := list.New(items, delegate, m.width, m.height)

	l.SetShowHelp(false)
	l.SetShowFilter(false)
	l.Styles.StatusBar.Foreground(gray)
	l.Styles.StatusBarFilterCount.Foreground(gray)
	l.Styles.StatusBarFilterCount.Foreground(gray)
	l.Styles.FilterPrompt.Foreground(yellow)

	return l
}

func (m Clide) PromptPath(n *node.CommandNode) (Clide, tea.Cmd) {
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
		state:  ClideStatePathSelect,
		list:   m.newlist(items),
	}, nil
}

type item struct {
	name, desc, value string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.value }

func (m Clide) PromptSelect() (Clide, tea.Cmd) {
	name := m.params[m.param].Name

	sibling := path.HasSibling(m.node.Path, name)

	if sibling == "" {
		return m.Error(fmt.Sprintf("Invalid parameter: No %s found in %s", name, path.Parent(m.node.Path)))
	}

	cmd := exec.Command(sibling)

	output, err := cmd.Output()

	if err != nil {
		return m.Error(fmt.Sprintf("Could not execute command %s", sibling))
	}

	options := strings.Split(string(output), "\n")

	if len(options) == 0 {
		return m.Error(fmt.Sprintf("%s yielded no options", name))
	}

	items := make([]list.Item, len(options))

	for i, choice := range options {
		if choice != "" {
			values := strings.Split(choice, ":")

            value := values[0]
            if len(values) >= 3 {
                value = values[2]
            }

			items[i] = list.Item(item{values[0], values[1], value})
		}
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
		list:   m.newlist(items),
	}, nil
}

func (m Clide) PromptInput() (Clide, tea.Cmd) {
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
	c.textinput.Placeholder = "..."
	c.textinput.Prompt = ""
	c.textinput.PromptStyle.BorderBottom(true)

	c.textinput.Focus()

	return c, nil
}
