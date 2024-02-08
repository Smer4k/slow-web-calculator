package main

import (
	"net/http"

	"github.com/Smer4k/slow-web-calculator/internal/agent"
)

func main() {
	agent := agent.NewAgent("http://localhost:8080/")
	agent.InitAgent()
	http.ListenAndServe(":9090", nil)
}