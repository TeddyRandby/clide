package model

import (
	"fmt"

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
			Background(purple).
			Foreground(bg).
			Margin(0, 1).
			Padding(0, 1)
	}()
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.
			NewStyle().
			BorderStyle(b).
			Background(bg).
			Foreground(purple).
			Padding(0, 1)
	}()
	stepStyle = func() lipgloss.Style {
		return lipgloss.
			NewStyle().
			ColorWhitespace(false).
			Foreground(white).
			Background(bg).
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
            Background(bg).
            Padding(0, 1)
    }()
    doneStyle = func() lipgloss.Style {
        return lipgloss.
            NewStyle().
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
		return fmt.Sprintf("%s\n%s", headerView, m.list.View())

	case ClideStatePromptSelect:
		m.list.SetSize(m.width, m.height-verticalSpace)
        m.list.SetShowTitle(false)
		return fmt.Sprintf("%s\n%s", headerView, m.list.View())

	case ClideStatePromptInput:
		m.textinput.Width = m.width
		header := promptStyle.Render(m.params[m.param].Name)
		return fmt.Sprintf("%s\n%s\n %s", headerView, header, m.textinput.View())

	case ClideStateError:
		return fmt.Errorf("Clide Error: %s.\n\nPress any key to exit.", m.error).Error()
	}

	return "Internal Error: Unknown state"
}
