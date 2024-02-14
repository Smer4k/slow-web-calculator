package main

import (
	"flag"
	"net/http"

	"github.com/Smer4k/slow-web-calculator/internal/agent"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "9090", "Порт запускаемого сервера (агент)")
	flag.Parse()
}

func main() {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	agent := agent.NewAgent("http://localhost:8080/", port)
	agent.InitAgent()
	http.ListenAndServe(":"+port, nil)
}
