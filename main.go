package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	windowWidth  = 160
	windowHeight = 40
)

type model struct {
	textInput    textinput.Model
	output       []string
	history      []string
	historyIndex int
	cmd          *exec.Cmd
	cmdOutput    chan string
	cursor       string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "jchTerm> "
	ti.Focus()
	ti.Width = windowWidth - 2 // Account for prompt
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	return model{
		textInput:    ti,
		output:       []string{"Welcome to jchTerm!"},
		history:      []string{},
		historyIndex: -1,
		cmdOutput:    make(chan string),
		cursor:       "█", // Block cursor
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			if m.cmd != nil {
				m.cmd.Process.Signal(syscall.SIGINT)
				return m, nil
			}
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

	case string:
		m.output = append(m.output, strings.Split(msg, "\n")...)
		if len(m.output) > windowHeight-1 {
			m.output = m.output[len(m.output)-(windowHeight-1):]
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	outputStyle := lipgloss.NewStyle().
		Width(windowWidth).
		Height(windowHeight - 1) // Reserve one line for input

	outputText := strings.Join(m.output, "\n")
	if len(outputText) > (windowWidth * (windowHeight - 1)) {
		outputText = outputText[len(outputText)-(windowWidth*(windowHeight-1)):]
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		outputStyle.Render(outputText),
		m.textInput.View(),
	)
}

func (m *model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		switch command {
		case "clear":
			m.output = []string{}
			return nil
		case "rebuild":
			cmd := exec.Command("go", "build", "-o", "jchterm_new")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Sprintf("Build failed: %v\n%s", err, output)
			}
			return "jchTerm rebuilt. Restarting...\n" + m.restart()
		case "edit":
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nano" // default to nano if EDITOR is not set
			}
			cmd := exec.Command(editor, "main.go")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return fmt.Sprintf("Error opening editor: %v", err)
			}
			return "File edited. Use 'rebuild' to apply changes."
		default:
			m.cmd = exec.Command("sh", "-c", command)
			output, err := m.cmd.CombinedOutput()
			m.cmd = nil
			if err != nil {
				return fmt.Sprintf("Command exited with error: %v\n%s", err, output)
			}
			return string(output)
		}
	}
}

func (m model) restart() string {
	binary, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("Error getting executable path: %v", err)
	}
	args := os.Args
	env := os.Environ()
	err = syscall.Exec(binary, args, env)
	if err != nil {
		return fmt.Sprintf("Error restarting: %v", err)
	}
	return "" // This line will never be reached if restart is successful
}

func main() {
	// Set up the terminal
	output := termenv.NewOutput(os.Stdout)
	output.SetWindowTitle("jchTerm")

	// Set the font (note: this may not work in all terminals)
	fmt.Print("\x1b]50;SetProfile=Operator Mono Medium 12\x07")

	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
