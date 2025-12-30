package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dlcuy22/endmi/core"
	"github.com/dlcuy22/endmi/extensions"
)

type step int

const (
	stepProjectName step = iota
	stepTemplate
	stepCreating
	stepDone
)

type outputMsg struct {
	line string
}

type doneMsg struct {
	err error
}

type model struct {
	step        step
	projectName string
	cursor      int
	templates   []extensions.Template
	input       string
	err         error
	output      []string
	app         *core.App
}

func initialModel(app *core.App, templates []extensions.Template, projectName string) model {
	startStep := stepProjectName
	if projectName != "" {
		startStep = stepTemplate
	}

	return model{
		step:        startStep,
		projectName: projectName,
		templates:   templates,
		cursor:      0,
		input:       projectName,
		output:      []string{},
		app:         app,
	}
}

// NewProgram wires a Bubble Tea program for the CLI.
func NewProgram(app *core.App, templates []extensions.Template, projectName string) *tea.Program {
	m := initialModel(app, templates, projectName)
	p := tea.NewProgram(&m)

	app.Output = func(line string) {
		if p != nil {
			p.Send(outputMsg{line: line})
		}
	}

	return p
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == stepCreating {
				return m, nil
			}
			return m, tea.Quit

		case "enter":
			switch m.step {
			case stepProjectName:
				if m.input != "" {
					m.projectName = m.input
					m.step = stepTemplate
				}
			case stepTemplate:
				m.step = stepCreating
				return m, m.createProject()
			}

		case "up":
			if m.step == stepTemplate && m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.step == stepTemplate && m.cursor < len(m.templates)-1 {
				m.cursor++
			}

		case "backspace":
			if m.step == stepProjectName && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			if m.step == stepProjectName && len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}

	case outputMsg:
		m.output = append(m.output, msg.line)
		return m, nil

	case doneMsg:
		m.err = msg.err
		m.step = stepDone
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) View() string {
	var b strings.Builder

	b.WriteString("Endmi - Golang Project Manager\n\n")

	switch m.step {
	case stepProjectName:
		b.WriteString("Enter project name:\n")
		b.WriteString(fmt.Sprintf("> %s█\n\n", m.input))
		b.WriteString("Press Enter to continue")

	case stepTemplate:
		b.WriteString(fmt.Sprintf("Project: %s\n\n", m.projectName))
		b.WriteString("Select template:\n\n")
		for i, t := range m.templates {
			line := fmt.Sprintf("%s — %s", t.Name(), t.Description())
			if m.cursor == i {
				b.WriteString(fmt.Sprintf("\033[48;5;240m\033[97m > %s \033[0m\n", line))
			} else {
				b.WriteString(fmt.Sprintf("   %s\n", line))
			}
		}
		b.WriteString("\nUse ↑/↓ to navigate, Enter to select")

	case stepCreating:
		selected := m.templates[m.cursor]
		b.WriteString(fmt.Sprintf("Creating project '%s' with %s...\n\n", m.projectName, selected.Name()))
		b.WriteString("╭─ Output ─────────────────────────────────────╮\n")
		for _, line := range m.output {
			b.WriteString(fmt.Sprintf("│ \033[90m%s\033[0m\n", line))
		}
		b.WriteString("╰──────────────────────────────────────────────╯\n")

	case stepDone:
		if m.err != nil {
			b.WriteString(fmt.Sprintf("❌ Error: %v\n", m.err))
		} else {
			b.WriteString(fmt.Sprintf("✅ Project '%s' created successfully!\n\n", m.projectName))
			b.WriteString(fmt.Sprintf("cd %s && go run .\n", m.projectName))
		}
	}

	if m.step != stepDone {
		b.WriteString("\n\nPress ctrl+c or q to quit")
	}

	return b.String()
}

func (m *model) createProject() tea.Cmd {
	return func() tea.Msg {
		tmpl := m.templates[m.cursor]
		if err := m.app.CreateProject(tmpl, m.projectName); err != nil {
			return doneMsg{err: err}
		}
		return doneMsg{err: nil}
	}
}
