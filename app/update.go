package model

import (
	"fmt"
	"strings"

	"github.com/TeddyRandby/clide/node"
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

func (m Clide) SelectPath(name string) (Clide, tea.Cmd) {
	for _, child := range m.node.Children {

		if strings.HasPrefix(child.Name, name) || child.Shortcut == name {
			switch child.Type {
			case node.NodeTypeCommand:
				return m.Command(&child)

			case node.NodeTypeModule:
				return m.PromptPath(&child)
			}
		}
	}

	return m.Error(fmt.Sprintf("No command or module '%s'", name))
}

func (m Clide) updateInput(msg tea.Msg) (Clide, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Next):
			value := m.textarea.Value()
			return m.SetAndPrompParameter(value)
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack()
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Root):
			return m.Root()
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	return m, cmd
}

func (m Clide) updatePathSelect(msg tea.Msg) (Clide, tea.Cmd) {
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
				return m.Root()
			}
		case key.Matches(msg, m.keymap.Root):
			return m.Root()

		case key.Matches(msg, m.keymap.VimNext):
			fallthrough
		case key.Matches(msg, m.keymap.Next):
			if !m.list.SettingFilter() {
				return m.SelectPath(m.list.SelectedItem().FilterValue())
			}
		case key.Matches(msg, m.keymap.VimPrev):
			if !m.list.SettingFilter() {
				return m.Backtrack()
			}
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Clide) updateError(msg tea.Msg) (Clide, tea.Cmd) {
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
			return m.Root()

		case key.Matches(msg, m.keymap.VimPrev):
			fallthrough
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack()
		}
	}

	return m, nil
}

func (m Clide) updateSelect(msg tea.Msg) (Clide, tea.Cmd) {
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
				return m.Root()
			}
		case key.Matches(msg, m.keymap.Root):
			return m.Root()

		case key.Matches(msg, m.keymap.VimNext):
			fallthrough
		case key.Matches(msg, m.keymap.Next):
			if !m.list.SettingFilter() && m.Param().HasValue() {
        return m.PromptParameter()
			}

		case key.Matches(msg, m.keymap.VimPrev):
			if !m.list.SettingFilter() {
				return m.Backtrack()
			}
		case key.Matches(msg, m.keymap.Prev):
			return m.Backtrack()

    case key.Matches(msg, m.keymap.Select):
      m.Set(m.list.SelectedItem().FilterValue())
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

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
		return m.updateError(msg)

	case ClideStatePathSelect:
		return m.updatePathSelect(msg)

	case ClideStatePromptSelect:
		return m.updateSelect(msg)

	case ClideStatePromptInput:
		return m.updateInput(msg)
	}

	panic("unreachable")
}
