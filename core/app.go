// app.go
// Core application logic for project creation workflow.
//
// Types:
//   - OutputHandler: callback for streaming command output lines
//   - App: owns the project creation workflow with Debug mode support
//
// Functions:
//   - CreateProject: scaffolds a project using the provided template
//   - runCommandWithOutput: executes a command and streams output
//   - streamOutput: reads from a reader and sends lines to OutputHandler
//   - runShellCommand: runs a shell command string (cross-platform)
package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dlcuy22/endmi/extensions"
)

type OutputHandler func(line string)

type App struct {
	Output OutputHandler
	Debug  bool
}

func (a App) CreateProject(t extensions.Template, projectName string) error {
	// Check if the template uses an external init command
	if initCmd := t.InitCommand(); initCmd != "" {
		// Replace {{name}} placeholder with actual project name
		initCmd = strings.ReplaceAll(initCmd, "{{name}}", projectName)
		if err := a.runShellCommand(initCmd, "."); err != nil {
			return err
		}
		return nil
	}

	// Standard project creation (generate files manually)
	projectPath := projectName

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	baseDir := filepath.Join(projectPath, t.RootDir())
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	if err := a.runCommandWithOutput("go", projectPath, "mod", "init", projectName); err != nil {
		return err
	}

	for rel, content := range t.Files(projectName) {
		fullPath := filepath.Join(baseDir, rel)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	for _, dep := range t.Dependencies() {
		if err := a.runCommandWithOutput("go", projectPath, "get", dep); err != nil {
			return err
		}
	}

	if err := a.runCommandWithOutput("go", projectPath, "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

func (a App) runCommandWithOutput(name string, dir string, args ...string) error {
	if a.Debug {
		fmt.Printf("[DEBUG] Running: %s %s (in %s)\n", name, strings.Join(args, " "), dir)
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go a.streamOutput(stdout)
	go a.streamOutput(stderr)

	return cmd.Wait()
}

func (a App) streamOutput(r io.Reader) {
	if a.Output == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		a.Output(scanner.Text())
	}
}

func (a App) runShellCommand(command string, dir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Dir = dir

	if a.Debug {
		fmt.Printf("[DEBUG] Shell command: %s\n", command)
		fmt.Printf("[DEBUG] Working directory: %s\n", dir)
		fmt.Printf("[DEBUG] Full exec: %s %v\n", cmd.Path, cmd.Args)
	}

	// Use CombinedOutput to capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if a.Debug && len(output) > 0 {
		fmt.Printf("[DEBUG] Command output:\n%s\n", string(output))
	}
	if err != nil {
		if len(output) > 0 {
			fmt.Printf("Command output:\n%s\n", string(output))
		}
		return err
	}

	return nil
}
