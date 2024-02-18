package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	if strings.Contains(port, ":") {
		panic("Один из параметров содержит \":\". В параметрах запуска не нужно указывать двоеточие")
	}
	re := regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(port) {
		panic("Порт должен содержать только цифры")
	}
	if num, _ := strconv.Atoi(port); num > 65535 {
		panic("Порт не должен быть больше 65535")
	}
	port = ":" + port
	orchestratorPort := ":8080"

	resp, err := http.Get(fmt.Sprintf("http://localhost%s/", orchestratorPort)) // проверка адреса на действительность
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	agent := agent.NewAgent(orchestratorPort, port)
	agent.InitAgent()
	http.ListenAndServe(port, nil)
}
