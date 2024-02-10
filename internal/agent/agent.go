package agent

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Agent struct {
	Router         *mux.Router
	AddrMainServer string
}

func NewAgent(addr string) *Agent {
	return &Agent{
		Router:         mux.NewRouter(),
		AddrMainServer: addr,
	}
}

func (a *Agent) InitAgent() {
	a.Router.HandleFunc("/solvingExpression", a.solvingExpression).Methods(http.MethodPost)
	a.Router.HandleFunc("/solvingExpression", a.redirectToMainServer).Methods(http.MethodGet)
	a.Router.HandleFunc("/", a.redirectToMainServer).Methods(http.MethodGet)
	http.Handle("/", a.Router)
}
