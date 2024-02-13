package main

import (
	"flag"
	"net/http"

	"github.com/Smer4k/slow-web-calculator/internal/agent"
)

var (
	addrMainServer string
	port string
)

func init() {
	flag.StringVar(&port, "port", ":9090", "Порт запускаемого сервера (агент)")
	flag.StringVar(&addrMainServer, "server", "http://localhost:8080/", "URL главного сервера (оркестратор)")
	flag.Parse()
}

func main() {
	resp, err := http.Get(addrMainServer)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	agent := agent.NewAgent(addrMainServer, port)
	agent.InitAgent()
	http.ListenAndServe(port, nil)
}
