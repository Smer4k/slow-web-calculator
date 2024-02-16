package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Agent struct {
	Router         *mux.Router
	AddrMainServer string
	AddrAgent      string
	Status         datatypes.Status
	CurrentTask    datatypes.Task
}

func NewAgent(orchestratorPort, port string) *Agent {
	return &Agent{
		Router:         mux.NewRouter(),
		AddrMainServer: "http://localhost" + orchestratorPort + "/",
		AddrAgent:      "http://localhost" + port + "/",
	}
}

func (a *Agent) InitAgent() {
	a.Router.HandleFunc("/solvingExpression", a.solvingExpression).Methods(http.MethodPost)
	a.Router.HandleFunc("/solvingExpression", a.redirectToMainServer).Methods(http.MethodGet)
	a.Router.HandleFunc("/", a.redirectToMainServer).Methods(http.MethodGet)
	http.Handle("/", a.Router)
	a.AddAgentToMainServer()
	a.PingMainServer()
}

func (a *Agent) AddAgentToMainServer() {
	vals := url.Values{}
	vals.Add("server", a.AddrAgent)
	_, err := http.PostForm(a.AddrMainServer+"addServer", vals)
	if err != nil {
		panic(err)
	}
	a.Status = datatypes.Idle
}

func (a *Agent) GetNewTask() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sgetExpression?%s=%s", a.AddrMainServer, "agent", a.AddrAgent), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("agent", a.AddrAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var task datatypes.Task
	fmt.Println(resp.Body)
	if contentType := resp.Header.Get("Content-Type"); contentType == "application/json" {
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			fmt.Println("Ошибка при чтении тела ответа:", err)
			return
		}
	}
	fmt.Println(task)
}

func (a *Agent) PostAnswer() {

}

func (a *Agent) SolveExpression() {

}

func (a *Agent) PingMainServer() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			_, err := http.Get(a.AddrMainServer)
			if err != nil {
				a.Status = datatypes.Reconnect
				fmt.Println(err)
				continue
			} else {
				if a.Status == datatypes.Reconnect {
					a.AddAgentToMainServer()
				}
				if a.Status == datatypes.Idle {
					a.GetNewTask()
				}
			}
		}
	}()
}
