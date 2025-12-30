package extensions

func init() {
	RegisterTemplate(netHTTPTemplate{})
}

type netHTTPTemplate struct{}

func (netHTTPTemplate) Name() string        { return "net/http" }
func (netHTTPTemplate) Description() string { return "Standard library HTTP server" }
func (netHTTPTemplate) RootDir() string     { return "" }
func (netHTTPTemplate) Dependencies() []string {
	return nil
}
func (netHTTPTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		"main.go": `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string ` + "`json:\"message\"`" + `
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/ping", handlePing)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Hello from ` + projectName + `!",
	})
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "pong",
	})
}
`,
	}
}
