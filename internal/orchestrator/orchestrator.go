package orchestrator

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Data struct {
	Done bool
	List []string
}

type Orchestrator struct {
	Router *mux.Router
	Tmpl   *template.Template
}

var data *Data

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{
		Router: mux.NewRouter(),
		Tmpl:   template.Must(template.ParseGlob("../../templates/*.html")),
	}
	return o
}

func (o *Orchestrator) InitRoutes() {
	data = &Data{}
	o.Router.HandleFunc("/", o.handleGetIndex).Methods(http.MethodGet)
	o.Router.HandleFunc("/calculator", o.handleGetCalculator).Methods(http.MethodGet)
	o.Router.HandleFunc("/calculator", o.handlePostCalculator).Methods(http.MethodPost)
	o.Router.HandleFunc("/settings", o.handleGetSettings).Methods(http.MethodGet)
	o.Router.HandleFunc("/settings", o.handlePostSettings).Methods(http.MethodPost)
	o.Router.HandleFunc("/results", o.handleGetResult).Methods(http.MethodGet)
	o.Router.HandleFunc("/postResult", o.handlePostResult).Methods(http.MethodPost)
	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
}
