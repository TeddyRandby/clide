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
	error     string
	ready     bool
	width     int
	height    int
	textinput textinput.Model
	viewport  viewport.Model
	list      list.Model
	node      *node.CommandNode
	root      *node.CommandNode
	params    []node.CommandNodeParameters
	param     int
	args      *map[string]string
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
		m.list.Title = "Continue"
		return m.list.View()

	case ClideStatePromptSelect:
		m.list.Title = m.params[m.param].Name
		return m.list.View()

	case ClideStatePromptInput:
		return m.textinput.View()

	case ClideStateError:
		return fmt.Errorf("Clide Error: %s.\n\nPress any key to exit.", m.error).Error()

	case ClideStateDone:
		return fmt.Sprintf("Clide Executed: %s\n\n %s", m.node.Path, m.viewport.View())
	}

	return "Internal Error: Unknown state"
}

func New(args *map[string]string) Clide {
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
		args: args,
	}.PromptPath(root)
}

func (m Clide) Run() {
	p := tea.NewProgram(m)

	// Run returns the model as a tea.Model.
	_, err := p.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
