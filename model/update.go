package model

import (
	"clide/node"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Clide) Index(name string) int {
	for i, child := range m.node.Children {
		if child.Shortcut != "" && child.Shortcut == name {
			return i
		}

		if strings.HasPrefix(child.Name, name) {
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Next):
			value := m.textinput.Value()
			return m.SetAndPromptNextArgument(value), nil
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack(), nil
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Root):
			return m.Root(), nil
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	return m, cmd
}

func (m Clide) updatePathSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.VimQuit):
			if !m.list.SettingFilter() {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.VimRoot):
			if !m.list.SettingFilter() {
				return m.Root(), nil
			}
		case key.Matches(msg, m.keymap.Root):
			return m.Root(), nil

		case key.Matches(msg, m.keymap.VimNext):
			if !m.list.SettingFilter() {
				return m.SelectPath(m.list.Index()), nil
			}
		case key.Matches(msg, m.keymap.Next):
			return m.SelectPath(m.list.Index()), nil

		case key.Matches(msg, m.keymap.VimPrev):
			if !m.list.SettingFilter() {
				return m.Backtrack(), nil
			}
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack(), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Clide) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.VimQuit):
			if !m.list.SettingFilter() {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.VimRoot):
			if !m.list.SettingFilter() {
				return m.Root(), nil
			}
		case key.Matches(msg, m.keymap.Root):
			return m.Root(), nil

		case key.Matches(msg, m.keymap.VimNext):
			if !m.list.SettingFilter() {
				return m.SetAndPromptNextArgument(m.list.SelectedItem().FilterValue()), nil
			}
		case key.Matches(msg, m.keymap.Next):
			return m.SetAndPromptNextArgument(m.list.SelectedItem().FilterValue()), nil

		case key.Matches(msg, m.keymap.VimPrev):
			if !m.list.SettingFilter() {
				return m.Backtrack(), nil
			}
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack(), nil

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd

}

func (m Clide) updateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.VimQuit):
			fallthrough
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.VimRoot):
			fallthrough
		case key.Matches(msg, m.keymap.Root):
			return m.Root(), nil

		case key.Matches(msg, m.keymap.VimPrev):
			fallthrough
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack(), nil
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
