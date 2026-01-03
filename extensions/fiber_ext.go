package extensions

func init() {
	RegisterTemplate(fiberTemplate{})
}

type fiberTemplate struct{}

func (fiberTemplate) Name() string        { return "fiber" }
func (fiberTemplate) Description() string { return "fiber template" }
func (fiberTemplate) RootDir() string     { return "" }
func (fiberTemplate) Dependencies() []string {
	return []string{"github.com/gofiber/fiber/v2", "github.com/gofiber/template/html/v2"}
}
func (fiberTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		"main.go": `package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":   "Hello from ` + projectName + `",
			"Message": "Fiber is running!",
		})
	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
`,
		"views/index.html": `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
</head>
<body>
  <h1>{{.Title}}</h1>
  <p>{{.Message}}</p>
</body>
</html>
`,
	}
}
