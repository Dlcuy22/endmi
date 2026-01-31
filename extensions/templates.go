// templates.go
// Template interface definition and registry for project templates.
//
// Types:
//   - Template: interface describing a project template
//
// Functions:
//   - RegisterTemplate: adds a template to the builtin registry
//   - BuiltinTemplates: returns the default templates bundled with the app
package extensions

type Template interface {
	Name() string
	Description() string
	RootDir() string
	Files(projectName string) map[string]string
	Dependencies() []string
	InitCommand() string
	PreCreateDir() bool
}

var registry []Template

func RegisterTemplate(t Template) {
	registry = append(registry, t)
}

func BuiltinTemplates() []Template {
	return registry
}
