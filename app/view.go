package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

var (
	fg     = lipgloss.Color("#F8F8F2")
	bg     = lipgloss.Color("#282A36")
	red    = lipgloss.Color("#FF5555")
	orange = lipgloss.Color("#FFB86C")
	yellow = lipgloss.Color("#F1FA8C")
	green  = lipgloss.Color("#50FA7B")
	cyan   = lipgloss.Color("#8BE9FD")
	purple = lipgloss.Color("#BD93F9")
	pink   = lipgloss.Color("#FF79C6")
	white  = lipgloss.Color("#ABB2BF")
	gray   = lipgloss.Color("#6272A4")
	black  = lipgloss.Color("#191A21")
)

var clide_header = "c l i d e"

var (
	titleStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(purple).
			Margin(0, 0, 0, 2)
	}()
	promptStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(orange).
			Margin(1, 1, 0).
			Padding(0, 1)
	}()
	stepStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			ColorWhitespace(false).
			Foreground(white).
			Margin(0, 1).
			Padding(0, 1)
	}()
	sepStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(purple)
	}()
	helpStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(gray).
			Padding(0, 1)
	}()
	errorStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(red).
			Padding(1, 1).
			Margin(0, 1)
	}()
	spinnerStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(red).
			Padding(0, 1)
	}()
)

func (m Clide) headerView() string {
	var steps []string

	node := m.node

	if node == nil {
		return ""
	}

	for node.Parent != nil {
		steps = append(steps, sepStyle.Render("/"), stepStyle.Render(node.Name))
		node = node.Parent
	}

	steps = append(steps, sepStyle.Render(clide_header))

	slices.Reverse(steps)

	return lipgloss.JoinHorizontal(lipgloss.Center, steps...)
}

func (m Clide) promptView() string {
	var params []string

	for i := 0; i < m.param+1; i++ {
		var str string

		if m.params[i].Value != "" {
			str = stepStyle.Render(m.params[i].Name)
		} else {
			str = promptStyle.Render(m.params[i].Name)
		}

		params = append(params, str, sepStyle.Render("/"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, params...)
}

func (m Clide) helpView() string {
	m.help.Styles.ShortKey.Foreground(gray)
	m.help.Styles.ShortDesc.Foreground(gray)
	return helpStyle.Render(m.help.View(m.keymap))
}

func (m Clide) View() string {
	if !m.ready {
		return "Initializing..."
	}

	m.help.Width = m.width

	helpView := m.helpView()
	headerView := m.headerView()

	verticalSpace := lipgloss.Height(headerView) + lipgloss.Height(helpView)

	switch m.state {

	case ClideStateStart:
		fallthrough
	case ClideStateDone:
		return ""

	case ClideStatePathSelect:
		m.list.SetSize(m.width, m.height-verticalSpace-2)

		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerView,
			m.list.View(),
			helpView,
		)

	case ClideStatePromptSelect:
		m.list.SetSize(m.width, m.height-verticalSpace-1)

		return fmt.Sprintf(
			"%s\n\n%s\n%s",
			lipgloss.JoinHorizontal(lipgloss.Right, headerView, m.promptView()),
			m.list.View(),
			helpView,
		)

	case ClideStatePromptInput:
		m.textarea.SetWidth(m.width)

		spaceRemaining := m.height - verticalSpace - 2

		textareaHeight := min(spaceRemaining/2, 8)

		m.textarea.SetHeight(textareaHeight)

		padding := spaceRemaining - textareaHeight - 1

		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Right, headerView, m.promptView()),
			"\n",
			m.textarea.View(),
			strings.Repeat("\n", padding),
			helpView,
		)

	case ClideStateError:
		content := errorStyle.
			Copy().
			Height(m.height - verticalSpace - 2).
			Render(fmt.Sprintf("Clide Error: %s.", m.error))

		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerView,
			content,
			helpView)
	}

	panic("unreachable")
}
