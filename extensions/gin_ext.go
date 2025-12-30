package extensions

func init() {
	RegisterTemplate(ginTemplate{})
}

type ginTemplate struct{}

func (ginTemplate) Name() string        { return "gin" }
func (ginTemplate) Description() string { return "Gin Web Framework" }
func (ginTemplate) RootDir() string     { return "" }
func (ginTemplate) Dependencies() []string {
	return []string{"github.com/gin-gonic/gin"}
}
func (ginTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		"main.go": `package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from ` + projectName + `!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
`,
	}
}
