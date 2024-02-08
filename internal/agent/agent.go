package agent

import (
	"net/http"
	"net/url"

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

func (a *Agent) solvingExpression(w http.ResponseWriter, r *http.Request) {
	val := url.Values{}
	val.Add("testVal", r.FormValue("expression"))
	resp, err := http.PostForm(a.AddrMainServer + "postResult", val)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// Middleware чтобы пользователь не смог зайти на сервер агента
func (a *Agent) redirectToMainServer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.AddrMainServer, http.StatusSeeOther)
}
