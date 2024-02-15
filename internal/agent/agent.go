package agent

import (
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

func NewAgent(addr, port string) *Agent {
	return &Agent{
		Router:         mux.NewRouter(),
		AddrMainServer: addr,
		AddrAgent:      "http://localhost" + port + "/",
	}
}

func (a *Agent) InitAgent() {
	a.Router.HandleFunc("/solvingExpression", a.solvingExpression).Methods(http.MethodPost)
	a.Router.HandleFunc("/solvingExpression", a.redirectToMainServer).Methods(http.MethodGet)
	a.Router.HandleFunc("/", a.redirectToMainServer).Methods(http.MethodGet)
	http.Handle("/", a.Router)
	a.AddAgentToMainServer()
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
			if err == nil {
				a.Status = datatypes.Reconnect
				fmt.Println(err)
				continue
			}
			if a.Status == datatypes.Reconnect {
				a.AddAgentToMainServer()
			}
		}
	}()
}
