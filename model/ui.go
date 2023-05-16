package model

import (
	"clide/node"
	"fmt"

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
		ready:  true,
		root:   m.root,
		state:  ClideStateError,
		error:  err,
	}
}

func (m Clide) Done(message string) Clide {
	clide := Clide{
		ready:    m.ready,
		width:    m.width,
		height:   m.height,
		node:     m.node,
		root:     m.root,
		state:    ClideStateDone,
		viewport: viewport.New(m.width, m.height),
	}

	clide.viewport.SetContent(message)
	return clide
}

func (m Clide) PromptPath(n *node.CommandNode) Clide {
    m.node = n

    if (m.node == nil) {
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
		prompt: "Select a command:",
		state:  ClideStatePathSelect,
		list:   list.New(items, list.NewDefaultDelegate(), m.width, m.height),
	}

	return c
}

func (m Clide) PromptSelect(prompt string, choices []string) Clide {
	return Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		node:   m.node,
		root:   m.root,
		prompt: prompt,
		state:  ClideStatePromptSelect,
	}
}

func (m Clide) PromptInput(prompt string) Clide {
    c := Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		node:   m.node,
		root:   m.root,
		prompt: prompt,
		state:  ClideStatePromptInput,
        textinput: textinput.New(),
	}

    c.textinput.Width = m.width
    c.textinput.CharLimit = 50
    c.textinput.Placeholder = "Enter a value"

    c.textinput.Focus()

    return c
}
