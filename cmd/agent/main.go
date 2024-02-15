package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/Smer4k/slow-web-calculator/internal/agent"
)

var (
	port string
	mainServerPort string
)

func init() {
	flag.StringVar(&port, "port", "9090", "Порт запускаемого сервера (агент)")
	flag.StringVar(&mainServerPort, "mainserver", "8080", "Порт главного сервера (оркестра)")
	flag.Parse()
}

func main() {
	if strings.Contains(port, ":") || strings.Contains(mainServerPort, ":") {
		panic("Один из параметров содержит \":\". В параметрах запуска не нужно указывать двоеточие")
	}
	port = ":" + port
	mainServerPort = ":" + mainServerPort

	resp, err := http.Get(fmt.Sprintf("http://localhost%s/ping", mainServerPort)) // проверка адреса на действительность
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if val := resp.Header.Get("answer"); val != "pong" { // исправить
		panic("Это не оркестар")
	}
	agent := agent.NewAgent(fmt.Sprintf("http://localhost%s/", mainServerPort), port)
	agent.InitAgent()
	http.ListenAndServe(port, nil)
}
