package main

import (
	"net/http"
	"path/filepath"

	"github.com/Smer4k/slow-web-calculator/internal/orchestrator"
)

// запускает storage (главный сервер)
func main() { 
	dir, _ := filepath.Abs(".")
	file, _ := filepath.Glob(filepath.Join(dir, "/main.go"))
	if len(file) == 0 {
		panic("Сервер можно запустить только из папки /cmd/orchestrator")
	}
	o := orchestrator.NewOrchestrator()
	o.InitRoutes()
	port := ":8080"
	http.ListenAndServe(port, nil)
}
