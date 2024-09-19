package commands

import (
	"strings"
)

type Commander interface {
	Execute(args []string) (string, error)
}

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
