---
description: How to add a new project template extension
---

# Adding a new template

1) Create a file named `<name>_ext.go` in this directory.
2) Implement the `Template` interface from `templates.go`.
3) In the same file, call `RegisterTemplate(...)` inside an `init()` to register it.

Example skeleton:

```go
package extensions

func init() {
    RegisterTemplate(myTemplate{})
}

type myTemplate struct{}

func (myTemplate) Name() string        { return "my-template" }
func (myTemplate) Description() string { return "Describe what it builds" }
func (myTemplate) Dependencies() []string {
    return []string{"example.com/some/dep"}
}
func (myTemplate) Files(projectName string) map[string]string {
    return map[string]string{
        "main.go": ` + "`package main\n\nfunc main() {}\n`" + `,
    }
}
```

Notes:

- Use `Files` to return all files to write (key = relative path, value = content).
- Use `Dependencies` for any modules needed; they will be `go get`-ed and `go mod tidy` will run.
- The template name is what appears in the UI list.
