package model

import (
	"clide/node"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Clide) Index(name string) int {
    for i, child := range m.node.Children {
        if strings.Contains(child.Name, name) {
            return i
        }
    }

    return -1
}

func (m Clide) SelectPath(index int) Clide {
	n := m.node.Children[index]

	switch n.Type {
	case node.NodeTypeCommand:
		return m.Command(&n)

	case node.NodeTypeModule:
		return m.PromptPath(&n)
	}

	return m
}

func (m Clide) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textinput.Width = m.width

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := m.textinput.Value()
			return m.SetAndPromptNextArgument(value), nil
		case "esc":
			return m.Backtrack(), nil
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	return m, cmd
}

func (m Clide) updatePathSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.list.SetSize(m.width, m.height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", " ":
			m := m.SelectPath(m.list.Index())
			return m, nil
		case "backspace", "h":
			return m.Backtrack(), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Clide) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.list.SetSize(m.width, m.height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", " ":
			m := m.SetAndPromptNextArgument(m.list.SelectedItem().FilterValue())
			return m, nil
		case "backspace", "h":
			return m.Backtrack(), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd

}

func (m Clide) updateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.viewport.Width = m.width
	m.viewport.Height = m.height

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "r":
			return m.Root(), nil
		}

	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m Clide) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	switch m.state {

	case ClideStateError:
		return m, tea.Quit

	case ClideStateDone:
		return m.updateDone(msg)

	case ClideStatePathSelect:
		return m.updatePathSelect(msg)

	case ClideStatePromptSelect:
		return m.updateSelect(msg)

	case ClideStatePromptInput:
		return m.updateInput(msg)
	}

	return m, nil
}
