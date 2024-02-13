package agent

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Agent struct {
	Router        *mux.Router
	AddrMainServer string
	AddrAgent      string
}

func NewAgent(addr, port string) *Agent {
	return &Agent{
		Router:         mux.NewRouter(),
		AddrMainServer: addr,
		AddrAgent: "http://localhost:" + port + "/",
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
	_, err := http.PostForm(a.AddrMainServer + "addServer", vals)
	if err != nil {
		panic(err)
	}
}