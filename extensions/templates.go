package extensions

// Template describes a project template that can generate source files and
// declare dependencies required to build the project.
//
// Implementations live in their own files (e.g., gin_ext.go, nethttp_ext.go)
// and call RegisterTemplate in an init() function to be included.
type Template interface {
	// Name is the selector key shown to the user.
	Name() string
	// Description provides a human-friendly summary.
	Description() string
	// RootDir returns the relative directory (under the project root) where
	// template files should be placed. Return an empty string to use the project
	// root.
	RootDir() string
	// Files returns the set of files to write for the project. The map key is
	// the relative path (e.g., "main.go") and the value is the file content.
	Files(projectName string) map[string]string
	// Dependencies lists Go modules that should be installed with `go get`.
	Dependencies() []string
	// InitCommand returns an optional shell command to run for project
	// initialization (e.g., "wails init -n {{name}} -t react-ts"). The
	// placeholder {{name}} will be replaced with the actual project name.
	// Return an empty string if no external command is needed.
	InitCommand() string
	// PreCreateDir returns whether to create the project directory before
	// running InitCommand. Set to false if the external tool creates its own
	// directory (e.g., wails init). Defaults to true for standard templates.
	PreCreateDir() bool
}

var registry []Template

// RegisterTemplate adds a template to the builtin registry. Call this from
// an init() inside each template file.
func RegisterTemplate(t Template) {
	registry = append(registry, t)
}

// BuiltinTemplates returns the default templates bundled with the app.
func BuiltinTemplates() []Template {
	return registry
}
