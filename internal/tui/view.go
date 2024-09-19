package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jchackert/jchTerm/internal/config"
)

func (m Model) View() string {
	outputStyle := lipgloss.NewStyle().
		Width(config.WindowWidth).
		Height(config.WindowHeight - 1)

	outputText := strings.Join(m.output, "\n")
	if len(outputText) > (config.WindowWidth * (config.WindowHeight - 1)) {
		outputText = outputText[len(outputText)-(config.WindowWidth*(config.WindowHeight-1)):]
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		outputStyle.Render(outputText),
		m.textInput.View(),
	)
}

func Run() error {
	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}

4. internal/commands/commands.go
package commands

import (
	"strings"
)

func ExecuteCommand(command string) (string, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return "", nil
	}

	switch args[0] {
	case "clear":
		return executeClear()
	case "rebuild":
		return executeRebuild()
	case "edit":
		return executeEdit()
	case "ask":
		return executeAsk(args[1:])
	default:
		return executeShell(args)
	}
}
