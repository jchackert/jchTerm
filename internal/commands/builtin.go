package commands

import (
	"os"
	"os/exec"
)

func executeClear() (string, error) {
	return "", nil
}

func executeRebuild() (string, error) {
	cmd := exec.Command("go", "build", "-o", "jchterm_new", "cmd/jchterm/main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return "jchTerm rebuilt. Please restart the application.", nil
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
		return "", err
	}
	return "File edited. Use 'rebuild' to apply changes.", nil
}

func executeShell(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
