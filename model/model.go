package model

import (
	"clide/node"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ClideStateStart        = 0
	ClideStatePathSelect   = 1
	ClideStatePromptSelect = 2
	ClideStatePromptInput  = 3
	ClideStateDone         = 4
	ClideStateError        = 5
)

type Clide struct {
	state     int
	prompt    string
	error     string
	ready     bool
	width     int
	height    int
	textinput textinput.Model
	viewport  viewport.Model
	list      list.Model
	node      *node.CommandNode
	root      *node.CommandNode
}

func (m Clide) Init() tea.Cmd {
	return nil
}

func (m Clide) View() string {
	if !m.ready {
		return "Initializing..."
	}

	switch m.state {

	case ClideStatePathSelect:
		fallthrough
	case ClideStatePromptSelect:
		m.list.Title = m.prompt
		return m.list.View()

    case ClideStatePromptInput:
        m.textinput.Prompt = m.prompt
        return m.textinput.View()

	case ClideStateError:
		return fmt.Errorf("Clide Error: %s.\n\nPress any key to exit.", m.error).Error()

	case ClideStateDone:
		return fmt.Sprintf("Clide Executed: %s\n\n %s", m.node.Path, m.viewport.View())
	}

	return "Internal Error: Unknown state"
}

func New() Clide {
	root, err := node.Root()

	if err != nil {
		return Clide{}.Error(err.Error())
	}

    if root == nil {
		return Clide{}.Error("No clide project found.")
    }

	return Clide{
		root: root,
		node: root,
	}
}

func (m Clide) Run() {
	c := m.PromptPath(m.node)

	p := tea.NewProgram(c)

	// Run returns the model as a tea.Model.
	_, err := p.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
