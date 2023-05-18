package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	black  = lipgloss.Color("#191A21")
)

var (
	promptStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(bg).
			Background(purple).
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
            Foreground(white).
            Padding(0, 1)
    }()
    errorStyle= func() lipgloss.Style {
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
	var renderedSteps []string

	node := m.node

    if node == nil {
        return ""
    }

	for node.Parent != nil {
		renderedSteps = append(renderedSteps, sepStyle.Render("/"), stepStyle.Render(node.Name))
		node = node.Parent
	}

	if len(renderedSteps) > 0 {
		renderedSteps = append(renderedSteps, sepStyle.Copy().MarginLeft(1).Render("󱞩"))
	}

	reversed := make([]string, len(renderedSteps))
	for i := 0; i < len(renderedSteps); i++ {
		reversed[i] = renderedSteps[len(renderedSteps)-i-1]
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, reversed...)
}

func (m Clide) helpView() string {
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

	case ClideStatePathSelect:
		m.list.SetSize(m.width, m.height-verticalSpace)
        m.list.SetShowTitle(false)
		return fmt.Sprintf("%s\n%s\n%s", headerView, m.list.View(), helpView)

	case ClideStatePromptSelect:
		m.list.SetSize(m.width, m.height-verticalSpace)
        m.list.SetShowTitle(false)
		return fmt.Sprintf("%s\n%s\n%s", headerView, m.list.View(), helpView)

	case ClideStatePromptInput:
		m.textinput.Width = m.width
		header := promptStyle.Render(m.params[m.param].Name)

        inputView := m.textinput.View()

        lines := m.height - verticalSpace - 1
        blank := strings.Repeat("\n", max(0, lines - lipgloss.Height(inputView)))
		return fmt.Sprintf("%s\n%s%s\n%s%s", headerView, header, inputView, blank, helpView)

	case ClideStateError:
        content := errorStyle.Copy().Height(m.height - verticalSpace).Render(fmt.Sprintf("Clide Error: %s.\n\n", m.error))
		return fmt.Sprintf("%s\n%s\n%s", headerView, content, helpView)
	}

	return "Internal Error: Unknown state"
}
