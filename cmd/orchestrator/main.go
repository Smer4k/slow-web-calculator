package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Smer4k/slow-web-calculator/internal/database"
	"github.com/Smer4k/slow-web-calculator/internal/orchestrator"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8080", "Порт запускаемого сервера (оркестор)")
	flag.Parse()
}

// запускает оркестр (главный сервер)
func main() {
	dir, _ := filepath.Abs(".")
	file, _ := filepath.Glob(filepath.Join(dir, "/main.go"))
	if len(file) == 0 {
		panic("Сервер можно запустить только из папки /cmd/orchestrator")
	}
	if strings.Contains(port, ":") {
		panic("Параметр запуска содержит \":\". В параметре запуска не нужно указывать двоеточие")
	}
	o := orchestrator.NewOrchestrator()

	database.InitDataBase()
	o.InitRoutes()

	port = ":" + port
	fmt.Printf("Сервер был успешно запущен и доступен по адресу \"http://localhost%s/\"\n", port)
	http.ListenAndServe(port, nil)
}
