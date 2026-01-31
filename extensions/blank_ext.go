package extensions

func init() {
	RegisterTemplate(blankTemplate{})
}

type blankTemplate struct{}

func (blankTemplate) Name() string        { return "blank" }
func (blankTemplate) Description() string { return "Empty Go project" }
func (blankTemplate) RootDir() string     { return "" }
func (blankTemplate) InitCommand() string { return "" }
func (blankTemplate) PreCreateDir() bool  { return true }
func (blankTemplate) Dependencies() []string {
	return nil
}
func (blankTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello from ` + projectName + `!")
}
`,
	}
}
