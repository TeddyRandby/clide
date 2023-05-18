package model

import (
	"clide/node"
	"fmt"
	"os"
	"syscall"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ClideStateStart        = 0
	ClideStatePathSelect   = 1
	ClideStatePromptSelect = 2
	ClideStatePromptInput  = 3
	ClideStateError        = 5
)

type KeyMap struct {
	Next    key.Binding
	Prev    key.Binding
	Root    key.Binding
	Quit    key.Binding
	VimNext key.Binding
	VimPrev key.Binding
	VimRoot key.Binding
	VimQuit key.Binding
}

var DefaultKeyMap = KeyMap{
	Next: key.NewBinding(
		key.WithKeys("right", "enter"),
		key.WithHelp("→/enter", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "esc"),
		key.WithHelp("←/esc", "previous"),
	),
	Root: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "root"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	VimNext: key.NewBinding(
		key.WithKeys("l", " "),
		key.WithHelp("l", "next"),
	),
	VimPrev: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "previous"),
	),
	VimRoot: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "root"),
	),
	VimQuit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Root, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Root, k.Quit},
		{k.VimNext, k.VimPrev},
	}
}

type Clide struct {
	state     int
	error     string
	ready     bool
	width     int
	height    int
	help      help.Model
	textinput textinput.Model
	list      list.Model
	spinner   spinner.Model
	node      *node.CommandNode
	root      *node.CommandNode
	params    []node.CommandNodeParameters
	param     int
	args      map[string]string
	keymap    KeyMap
}

func (m Clide) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Clide) Leaves() []node.CommandNode {
	return m.node.Leaves()
}

func New(args map[string]string) Clide {
	root, err := node.Root()

	if err != nil {
		m, _ := Clide{}.Error(err.Error())
		return m.(Clide)
	}

	if root == nil {
		m, _ := Clide{}.Error("No project found")
		return m.(Clide)
	}

	m, _ := Clide{
		root:   root,
		node:   root,
		args:   args,
		keymap: DefaultKeyMap,
		help:   help.New(),
	}.PromptPath(root)

	return m.(Clide)
}

func (m Clide) Run() {
	c, err := tea.NewProgram(m).Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m = c.(Clide)

	if m.state != ClideStateError {
		syscall.Exec(m.node.Path, []string{m.node.Name}, os.Environ())
	}
}
