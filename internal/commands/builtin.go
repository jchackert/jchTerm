// internal/commands/builtin.go

package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func executeClear() (string, error) {
	return "", nil
}

func executeRebuild() (string, error) {
	cmd := exec.Command("go", "build", "-o", "jchterm", "cmd/jchterm/main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, output)
	}
	return "jchterm rebuilt. Please restart the application.", nil
}

func executeEdit() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	cmd := exec.Command(editor, "cmd/jchterm/main.go")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error opening editor: %v", err)
	}
	return "File edited. Use 'rebuild' to apply changes.", nil
}

func executeShell(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %v\n%s", err, output)
	}
	return string(output), nil
}
