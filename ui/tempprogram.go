package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dlcuy22/endmi/core"
	"github.com/dlcuy22/endmi/extensions"
	"github.com/dlcuy22/endmi/utils"
)

type tempStep int

const (
	tempStepProjectName tempStep = iota
	tempStepTemplate
	tempStepCreating
	tempStepDone
	tempStepChoice
)

type tempModel struct {
	step        tempStep
	projectName string
	cursor      int
	templates   []extensions.Template
	input       string
	err         error
	output      []string
	tcm         *core.TempCodeManager
	resultPath  string
}

func initialTempModel(tcm *core.TempCodeManager, templates []extensions.Template) tempModel {
	return tempModel{
		step:      tempStepTemplate,
		templates: templates,
		cursor:    0,
		input:     "",
		output:    []string{},
		tcm:       tcm,
	}
}

// NewTempProgram creates a Bubble Tea program for temporary project creation
func NewTempProgram(tcm *core.TempCodeManager, templates []extensions.Template) *tea.Program {
	m := initialTempModel(tcm, templates)
	p := tea.NewProgram(&m)

	tcm.App.Output = func(line string) {
		if p != nil {
			p.Send(outputMsg{line: line})
		}
	}

	return p
}

func (m *tempModel) Init() tea.Cmd {
	return nil
}

func (m *tempModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == tempStepCreating {
				return m, nil
			}
			return m, tea.Quit

		case "enter":
			switch m.step {
			case tempStepProjectName:
				// Allow empty input for auto-generated name
				m.step = tempStepTemplate
			case tempStepTemplate:
				m.step = tempStepCreating
				return m, m.createTempProject()
			case tempStepChoice:
				if m.cursor == 0 {
					// Open terminal in temp folder
					return m, m.openTerminal()
				} else {
					// Exit
					return m, tea.Quit
				}
			case tempStepDone:
				return m, tea.Quit
			}

		case "up":
			if m.step == tempStepTemplate && m.cursor > 0 {
				m.cursor--
			}
			if m.step == tempStepChoice && m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.step == tempStepTemplate && m.cursor < len(m.templates)-1 {
				m.cursor++
			}
			if m.step == tempStepChoice && m.cursor < 1 {
				m.cursor++
			}

		case "backspace":
			if m.step == tempStepProjectName && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		case "tab":
			if m.step == tempStepTemplate {
				// Skip project name step
				m.step = tempStepProjectName
				m.input = m.projectName
			}

		default:
			if m.step == tempStepProjectName && len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}

	case outputMsg:
		m.output = append(m.output, msg.line)
		return m, nil

	case doneMsg:
		m.err = msg.err
		if msg.err == nil {
			// Success - show choice
			m.step = tempStepChoice
			m.cursor = 0 // Reset cursor for choice menu
		} else {
			// Error - go to done
			m.step = tempStepDone
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m *tempModel) View() string {
	var b strings.Builder

	b.WriteString("Endmi - Temporary Code Workspace\n\n")

	switch m.step {
	case tempStepProjectName:
		b.WriteString("Enter project name (leave empty for auto-generated):\n")
		b.WriteString(fmt.Sprintf("> %s█\n\n", m.input))
		b.WriteString("Press Enter to continue or Tab to go back")

	case tempStepTemplate:
		if m.projectName != "" && m.input != "" {
			b.WriteString(fmt.Sprintf("Project: %s\n\n", m.input))
		} else {
			b.WriteString("Project: [auto-generated]\n\n")
		}
		b.WriteString("Select template:\n\n")
		b.WriteString(RenderTemplateList(m.templates, m.cursor))
		b.WriteString("\nUse ↑/↓ to navigate, Enter to create, Tab to set project name")

	case tempStepCreating:
		selected := m.templates[m.cursor]
		projectDisplayName := m.projectName
		if projectDisplayName == "" {
			projectDisplayName = "[auto-generated]"
		}
		b.WriteString(fmt.Sprintf("Creating temporary project '%s' with %s...\n\n", projectDisplayName, selected.Name()))
		b.WriteString(RenderOutputBox(m.output))

	case tempStepChoice:
		b.WriteString("✅ Temporary project created successfully!\n\n")
		b.WriteString(fmt.Sprintf("📁 Location: %s\n\n", m.resultPath))
		b.WriteString("What would you like to do?\n\n")
		b.WriteString(RenderChoiceMenu(m.cursor, "Open terminal in temp folder", "Exit"))
		b.WriteString("\nUse ↑/↓ to navigate, Enter to select")

	case tempStepDone:
		if m.err != nil {
			b.WriteString(fmt.Sprintf("❌ Error: %v\n", m.err))
		}
	}

	if m.step != tempStepDone && m.step != tempStepChoice {
		b.WriteString("\n\nPress ctrl+c or q to quit")
	}

	return b.String()
}

func (m *tempModel) createTempProject() tea.Cmd {
	return func() tea.Msg {
		tmpl := m.templates[m.cursor]
		projectPath, err := m.tcm.CreateTempProject(tmpl, m.input)
		if err != nil {
			return doneMsg{err: err}
		}
		m.resultPath = projectPath
		return doneMsg{err: nil}
	}
}

func (m *tempModel) openTerminal() tea.Cmd {
	return func() tea.Msg {
		if err := utils.OpenTerminalInDirectory(m.resultPath); err != nil {
			return doneMsg{err: err}
		}
		return tea.Quit()
	}
}
