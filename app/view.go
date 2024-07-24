package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/TeddyRandby/clide/node"
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

	n := m.node

	if n == nil {
		return ""
	}

	for n.Parent != nil {
		steps = append(steps, sepStyle.Render("/"), stepStyle.Render(n.Name))
		n = n.Parent
	}

	steps = append(steps, sepStyle.Render(clide_header))

	slices.Reverse(steps)

	return lipgloss.JoinHorizontal(lipgloss.Center, steps...)
}

func (m Clide) promptView() string {
	var params []string

	for i := 0; i < m.param+1; i++ {
		var str string

		if i < m.param {
			str = stepStyle.Render(strings.Join(m.params[i].Value, ", "))
		} else {
			str = promptStyle.Render(m.preview())
		}

		params = append(params, str, sepStyle.Render("/"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, params...)
}

func (m Clide) preview() string {
	if len(m.params) == 0 {
		return ""
	}

	var str string

	switch m.Param().Type {
	case node.CommandNodeParamTypeInput:
		str = m.textarea.Value()
	case node.CommandNodeParamTypeSelect:
    str = strings.Join(m.Param().Value, ", ")
	}

  if str == "" {
    return m.Param().Name
  } else {
    return str
  }
}

func (m Clide) helpView() string {
	m.help.Styles.ShortKey.Foreground(gray)
	m.help.Styles.ShortDesc.Foreground(gray)
	return helpStyle.Render(m.help.View(m.keymap))
}

func (m Clide) View() string {
	m.help.Width = m.width

	helpView := m.helpView()
	headerView := m.headerView()

	verticalSpace := lipgloss.Height(headerView) + lipgloss.Height(helpView) + 1

	switch m.state {

	case ClideStateStart:
		fallthrough
	case ClideStateDone:
		return ""

	case ClideStatePathSelect:
		m.list.SetSize(m.width, m.height-verticalSpace)

		return lipgloss.JoinVertical(lipgloss.Left,
			headerView,
			m.list.View(),
			helpView,
		)

	case ClideStatePromptSelect:
		m.list.SetSize(m.width, m.height-verticalSpace)
		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Right, headerView, m.promptView()),
			m.list.View(),
			helpView,
		)

	case ClideStatePromptInput:
		m.textarea.SetWidth(m.width)

		spaceRemaining := m.height - verticalSpace

		m.textarea.SetHeight(spaceRemaining)

		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Right, headerView, m.promptView()),
			m.textarea.View(),
			helpView,
		)

	case ClideStateError:
		content := errorStyle.
			Copy().
			Height(m.height - verticalSpace).
			Render(fmt.Sprintf("Clide Error: %s.", m.error))

		return lipgloss.JoinVertical(lipgloss.Left,
			headerView,
			content,
			helpView)
	}

	panic("unreachable")
}
