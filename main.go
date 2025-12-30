package main

import (
	"fmt"
	"os"

	"github.com/dlcuy22/endmi/core"
	"github.com/dlcuy22/endmi/extensions"
	"github.com/dlcuy22/endmi/ui"
)

func showHelp() {
	fmt.Println("Endmi - Golang Project Manager")
	fmt.Println()
	fmt.Println("A modular CLI tool to bootstrap Golang projects with templates.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  endmi create [project-name]    Create a new Go project")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  endmi create                   Start interactive project creation")
	fmt.Println("  endmi create my-api            Create project named 'my-api'")
	fmt.Println()
	fmt.Println("Available templates:")
	for _, t := range extensions.BuiltinTemplates() {
		fmt.Printf("  - %-10s %s\n", t.Name(), t.Description())
	}
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	if os.Args[1] != "create" {
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println()
		showHelp()
		os.Exit(1)
	}

	var projectName string
	if len(os.Args) > 2 {
		projectName = os.Args[2]
	}

	app := &core.App{}
	templates := extensions.BuiltinTemplates()

	program := ui.NewProgram(app, templates, projectName)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
