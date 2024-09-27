package tui

import (
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jchackert/jchterm/internal/commands"
	"github.com/jchackert/jchterm/internal/config"
	"github.com/jchackert/jchterm/internal/logger"
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
			logger.Log("Executing command: %s", command)
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
	default:
		// Only log messages that are not related to cursor operations
		msgType := strings.ToLower(string(reflect.TypeOf(msg).Name()))
		if !strings.Contains(msgType, "blink") && !strings.Contains(msgType, "cursor") {
			logger.Log("Received message of type: %T", msg)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		logger.Log("Running command: %s", command)
		output, err := commands.ExecuteCommand(command)
		if err != nil {
			logger.Log("Command error: %v", err)
			return err.Error()
		}
		logger.Log("Command output: %s", output)
		return output
	}
}

func (m Model) View() string {
	var s strings.Builder
	for _, line := range m.output {
		s.WriteString(line + "\n")
	}
	s.WriteString(m.textInput.View())
	return s.String()
}

func Run() error {
	logger.Log("Entering tui.Run()")
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	// Set up a channel to receive OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle the signal
	go func() {
		<-c
		logger.Log("Received Ctrl+C, exiting...")
		p.Kill()
	}()

	logger.Log("Running program...")
	_, err := p.Run()
	if err != nil {
		logger.Log("Error running program: %v", err)
	}
	logger.Log("Exiting tui.Run()")
	return err
}
