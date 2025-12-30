package core

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dlcuy22/endmi/extensions"
)

// OutputHandler receives streaming lines from command execution.
type OutputHandler func(line string)

// App owns the project creation workflow.
type App struct {
	Output OutputHandler
}

// CreateProject scaffolds a project using the provided template.
func (a App) CreateProject(t extensions.Template, projectName string) error {
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
