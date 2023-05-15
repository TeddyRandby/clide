package model

import (
	. "clide/path"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
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
	state    int
	path     string
	prompt   string
	error    string
	ready    bool
	width    int
	height   int
	viewport viewport.Model
	list     list.Model
	nodes    []CommandNode
}

func (m Clide) Init() tea.Cmd {
	return nil
}

func (m Clide) SelectPath(index int) Clide {
	node := m.nodes[index]
	newPath := filepath.Join(m.path, node.Name)

	switch node.Type {
	case NodeTypeCommand:
		cmd := exec.Command(newPath)

		m = Clide{
			ready:    m.ready,
			width:    m.width,
			height:   m.height,
			state:    ClideStateDone,
			path:     newPath,
			viewport: viewport.New(m.width, m.height),
		}

		output, err := cmd.Output()

		if err != nil {
			return Clide{error: err.Error(), state: ClideStateError}
		}

		m.viewport.SetContent(string(output))
		return m

	case NodeTypeParameter:
		return Clide{path: newPath}

	case NodeTypeModule:
		m.path = newPath
		return m.PromptPath()
	}

	return m
}

func (m Clide) PromptPath() Clide {
	pathChoices, err := PathChoices(m.path)

	if err != nil {
		return Clide{error: err.Error(), state: ClideStateError, ready: true}
	}

	if len(pathChoices) == 0 {
		return Clide{error: fmt.Sprintf("No commands found in %s", m.path), state: ClideStateError, ready: true}
	}

	items := make([]list.Item, len(pathChoices))

	for i, choice := range pathChoices {
		items[i] = list.Item(choice)
	}

	c := Clide{
		ready:  m.ready,
		width:  m.width,
		height: m.height,
		path:   m.path,
		prompt: "Select a command:",
		state:  ClideStatePathSelect,
		nodes:  pathChoices,
		list:   list.New(items, list.NewDefaultDelegate(), m.width, m.height),
	}

	return c
}

func (m Clide) PromptSelect(prompt string, choices []string) (Clide, error) {
	return Clide{width: m.width, height: m.height, path: m.path, prompt: prompt, state: ClideStatePromptSelect, ready: m.ready}, nil
}

func (m Clide) PromptInput(prompt string) (Clide, error) {
	return Clide{width: m.width, height: m.height, path: m.path, prompt: prompt, state: ClideStatePromptInput, ready: m.ready}, nil
}

func (m Clide) UpdatePathSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(m.width, m.height)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m := m.SelectPath(m.list.Index())

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Clide) UpdateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = m.width
		m.viewport.Height = m.height
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
		return m.UpdateDone(msg)

	case ClideStatePathSelect:
		return m.UpdatePathSelect(msg)

	case ClideStatePromptInput:
	}

	return m, nil
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

	case ClideStateError:
		return fmt.Errorf("Clide Error: %s.\n\nPress any key to exit.", m.error).Error()

	case ClideStateDone:
		return fmt.Sprintf("Clide Executed: %s\n\n %s", m.path, m.viewport.View())
	}

	return "Internal Error: Unknown state"
}

func Run(root string) {
	c := Clide{path: root}

	c = c.PromptPath()

	p := tea.NewProgram(c)

	// Run returns the model as a tea.Model.
	_, err := p.Run()

	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
