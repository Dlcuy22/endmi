package ui

import (
	"fmt"

	"github.com/dlcuy22/endmi/extensions"
)

// RenderTemplateList renders a selectable list of templates with the cursor
func RenderTemplateList(templates []extensions.Template, cursor int) string {
	result := ""
	for i, t := range templates {
		line := fmt.Sprintf("%s — %s", t.Name(), t.Description())
		if cursor == i {
			result += fmt.Sprintf("\033[48;5;240m\033[97m > %s \033[0m\n", line)
		} else {
			result += fmt.Sprintf("   %s\n", line)
		}
	}
	return result
}

// RenderChoiceMenu renders a two-option choice menu with cursor
func RenderChoiceMenu(cursor int, option1, option2 string) string {
	result := ""
	if cursor == 0 {
		result += fmt.Sprintf("\033[48;5;240m\033[97m > %s \033[0m\n", option1)
	} else {
		result += fmt.Sprintf("   %s\n", option1)
	}

	if cursor == 1 {
		result += fmt.Sprintf("\033[48;5;240m\033[97m > %s \033[0m\n", option2)
	} else {
		result += fmt.Sprintf("   %s\n", option2)
	}
	return result
}

// RenderOutputBox renders a bordered output box with lines
func RenderOutputBox(lines []string) string {
	result := "╭─ Output ─────────────────────────────────────╮\n"
	for _, line := range lines {
		result += fmt.Sprintf("│ \033[90m%s\033[0m\n", line)
	}
	result += "╰──────────────────────────────────────────────╯\n"
	return result
}
