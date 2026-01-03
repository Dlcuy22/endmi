package utils

import (
	"fmt"

	"github.com/dlcuy22/endmi/extensions"
)

// FindTemplateByName searches for a template by name in the provided list
func FindTemplateByName(templates []extensions.Template, name string) (extensions.Template, error) {
	for _, t := range templates {
		if t.Name() == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("template '%s' not found", name)
}

// ListTemplateNames returns a formatted string of all available template names
func ListTemplateNames(templates []extensions.Template) string {
	result := "Available templates:\n"
	for _, t := range templates {
		result += fmt.Sprintf("  - %-10s %s\n", t.Name(), t.Description())
	}
	return result
}
