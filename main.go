package main

import (
	"fmt"
	"os"

	"log"

	"github.com/dlcuy22/endmi/core"
	"github.com/dlcuy22/endmi/extensions"
	"github.com/dlcuy22/endmi/ui"
	"github.com/dlcuy22/endmi/utils"
)

func showHelp() {
	fmt.Println("Endmi - Golang Project Manager")
	fmt.Println()
	fmt.Println("A modular CLI tool to bootstrap Golang projects with templates.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  endmi create [project-name] [flags]    Create a new Go project")
	fmt.Println("  endmi temp <command> [flags]           Manage temporary code workspace")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --template <name>                  Specify template (skip interactive selection)")
	fmt.Println("  -n, --name <name>                      Specify project name (for temp create)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  endmi create                           Start interactive project creation")
	fmt.Println("  endmi create my-api                    Create project named 'my-api'")
	fmt.Println("  endmi create my-api -t fiber           Create 'my-api' with fiber template")
	fmt.Println("  endmi temp create                      Create a new temporary project")
	fmt.Println("  endmi temp create -t gin               Create temp project with gin template")
	fmt.Println("  endmi temp create -t blank -n mytest   Create named temp project")
	fmt.Println("  endmi temp list                        List all temporary projects")
	fmt.Println("  endmi temp delete <name>               Delete a temporary project")
	fmt.Println("  endmi temp clean                       Remove all temporary projects")
	fmt.Println("  endmi temp promote <name> <path>       Move temp project to permanent location")
	fmt.Println()
	fmt.Println("Available templates:")
	for _, t := range extensions.BuiltinTemplates() {
		fmt.Printf("  - %-10s %s\n", t.Name(), t.Description())
	}
}

func main() {
	exists, err := utils.CheckConfigExists()
	if err != nil {
		log.Fatalf("failed to check config existence: %v", err)
	}

	if !exists {
		log.Println("config not found, creating default config...")

		if err := utils.WriteConfig(); err != nil {
			log.Fatalf("failed to create config file: %v", err)
		}

		log.Println("default config created successfully")
	} else {
		log.Println("config already exists, skipping creation")
	}

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		var projectName string
		var templateName string

		// Parse arguments and flags
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "--template" || arg == "-t" {
				if i+1 < len(os.Args) {
					templateName = os.Args[i+1]
					i++ // Skip next arg since it's the template value
				} else {
					fmt.Println("Error: --template/-t requires a template name")
					os.Exit(1)
				}
			} else if projectName == "" {
				projectName = arg
			}
		}

		app := &core.App{}
		templates := extensions.BuiltinTemplates()

		// If template is specified via flag, create project directly
		if templateName != "" {
			if projectName == "" {
				fmt.Println("Error: project name is required when using --template")
				fmt.Println("Usage: endmi create <project-name> --template <template-name>")
				os.Exit(1)
			}

			// Find the template
			selectedTemplate, err := utils.FindTemplateByName(templates, templateName)
			if err != nil {
				fmt.Printf("Error: %v\n\n", err)
				fmt.Print(utils.ListTemplateNames(templates))
				os.Exit(1)
			}

			// Create project directly
			fmt.Printf("Creating project '%s' with template '%s'...\n", projectName, templateName)
			if err := app.CreateProject(selectedTemplate, projectName); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n✅ Project '%s' created successfully!\n", projectName)
			fmt.Printf("   cd %s && go run .\n", projectName)
		} else {
			// Use interactive UI
			program := ui.NewProgram(app, templates, projectName)
			if _, err := program.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	case "help":
		showHelp()
		os.Exit(0)

	case "temp":
		if len(os.Args) < 3 {
			fmt.Println("Error: temp command requires a subcommand")
			fmt.Println()
			fmt.Println("Available subcommands:")
			fmt.Println("  create              Create a new temporary project")
			fmt.Println("  list                List all temporary projects")
			fmt.Println("  delete <name>       Delete a temporary project")
			fmt.Println("  clean               Remove all temporary projects")
			fmt.Println("  promote <name> <path> Move temp project to permanent location")
			os.Exit(1)
		}

		subcommand := os.Args[2]
		app := &core.App{}
		tcm := &core.TempCodeManager{App: app}

		switch subcommand {
		case "create":
			templates := extensions.BuiltinTemplates()
			if len(templates) == 0 {
				fmt.Println("Error: No templates available")
				os.Exit(1)
			}

			var templateName string
			var projectName string

			// Parse flags for temp create
			for i := 3; i < len(os.Args); i++ {
				arg := os.Args[i]
				if arg == "--template" || arg == "-t" {
					if i+1 < len(os.Args) {
						templateName = os.Args[i+1]
						i++
					}
				} else if arg == "--name" || arg == "-n" {
					if i+1 < len(os.Args) {
						projectName = os.Args[i+1]
						i++
					}
				}
			}

			// If template is specified via flag, create directly
			if templateName != "" {
				selectedTemplate, err := utils.FindTemplateByName(templates, templateName)
				if err != nil {
					fmt.Printf("Error: %v\n\n", err)
					fmt.Print(utils.ListTemplateNames(templates))
					os.Exit(1)
				}

				fmt.Printf("Creating temporary project with template '%s'...\n", templateName)
				projectPath, err := tcm.CreateTempProject(selectedTemplate, projectName)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("\n✅ Temporary project created successfully!\n")
				fmt.Printf("📁 Location: %s\n\n", projectPath)
				fmt.Println("ℹ️  This is a temporary workspace. Changes won't be tracked.")
				fmt.Println("   Use 'endmi temp promote <name> <path>' to make it permanent.")
			} else {
				// Use interactive UI
				program := ui.NewTempProgram(tcm, templates)
				if _, err := program.Run(); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}

		case "list":
			projects, err := tcm.ListTempProjects()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if len(projects) == 0 {
				fmt.Println("No temporary projects found.")
				os.Exit(0)
			}

			fmt.Println("Temporary Projects:")
			fmt.Println()
			for _, p := range projects {
				fmt.Printf("  Name:     %s\n", p.Name)
				fmt.Printf("  Template: %s\n", p.Template)
				fmt.Printf("  Created:  %s\n", p.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("  Path:     %s\n", p.Path)
				fmt.Println()
			}

		case "delete":
			if len(os.Args) < 4 {
				fmt.Println("Error: delete requires a project name")
				fmt.Println("Usage: endmi temp delete <name>")
				os.Exit(1)
			}

			projectName := os.Args[3]
			fmt.Printf("Deleting temporary project '%s'...\n", projectName)

			if err := tcm.DeleteTempProject(projectName); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Temporary project '%s' deleted successfully\n", projectName)

		case "clean":
			fmt.Print("Are you sure you want to delete ALL temporary projects? (y/N): ")
			var confirm string
			fmt.Scanln(&confirm)

			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled.")
				os.Exit(0)
			}

			fmt.Println("Cleaning all temporary projects...")
			if err := tcm.CleanAll(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✓ All temporary projects removed successfully")

		case "promote":
			if len(os.Args) < 5 {
				fmt.Println("Error: promote requires a project name and target path")
				fmt.Println("Usage: endmi temp promote <name> <path>")
				os.Exit(1)
			}

			projectName := os.Args[3]
			targetPath := os.Args[4]

			fmt.Printf("Promoting temporary project '%s' to '%s'...\n", projectName, targetPath)
			if err := tcm.PromoteTempProject(projectName, targetPath); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ Project promoted successfully to: %s\n", targetPath)
			fmt.Println("  The project is now permanent and no longer in the temp workspace.")

		default:
			fmt.Printf("Unknown temp subcommand: %s\n", subcommand)
			fmt.Println()
			fmt.Println("Available subcommands:")
			fmt.Println("  create              Create a new temporary project")
			fmt.Println("  list                List all temporary projects")
			fmt.Println("  delete <name>       Delete a temporary project")
			fmt.Println("  clean               Remove all temporary projects")
			fmt.Println("  promote <name> <path> Move temp project to permanent location")
			os.Exit(1)
		}
		os.Exit(0)

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println()
		showHelp()
		os.Exit(1)
	}
}
