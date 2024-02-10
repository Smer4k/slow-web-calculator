package orchestrator

import (
	"net/http"
	"net/url"
)

func (o *Orchestrator) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Calculator.html
func (o *Orchestrator) handleGetCalculator(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "calculator.html", o.Data)
}

func (o *Orchestrator) handlePostCalculator(w http.ResponseWriter, r *http.Request) {
	vals := url.Values{}
	expression := r.FormValue("expression")
	vals.Add("expression", expression)
	resp, err := http.PostForm("http://localhost:9090/solvingExpression", vals)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	http.Redirect(w, r, "/calculator", http.StatusSeeOther)
}

// Settings.html
func (o *Orchestrator) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "settings.html", o.Data)
}

func (o *Orchestrator) handlePostSettings(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

// Results.html
func (o *Orchestrator) handleGetResult(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "results.html", o.Data)
}

func (o *Orchestrator) handlePostResult(w http.ResponseWriter, r *http.Request) {

}