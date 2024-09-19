package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jchackert/jchterm/internal/commands"
	"github.com/jchackert/jchterm/internal/config"
)
type Model struct {
	textInput    textinput.Model
	output       []string
	history      []string
	historyIndex int
	cursor       string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "jchTerm> "
	ti.Focus()
	ti.Width = config.WindowWidth - 2
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	return Model{
		textInput:    ti,
		output:       []string{"Welcome to jchTerm!"},
		history:      []string{},
		historyIndex: -1,
		cursor:       "█",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			command := m.textInput.Value()
			m.history = append(m.history, command)
			m.historyIndex = len(m.history)
			m.output = append(m.output, "$ "+command)
			m.textInput.SetValue("")
			return m, m.runCommand(command)
		case tea.KeyUp:
			if m.historyIndex > 0 {
				m.historyIndex--
				m.textInput.SetValue(m.history[m.historyIndex])
			}
		case tea.KeyDown:
			if m.historyIndex < len(m.history)-1 {
				m.historyIndex++
				m.textInput.SetValue(m.history[m.historyIndex])
			} else {
				m.historyIndex = len(m.history)
				m.textInput.SetValue("")
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		output, err := commands.ExecuteCommand(command)
		if err != nil {
			return err.Error()
		}
		return output
	}
}
