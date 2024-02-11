package orchestrator

import (
	"fmt"
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
	defer http.Redirect(w, r, "/calculator", http.StatusSeeOther)
	vals := url.Values{}
	expression := r.FormValue("expression")
	ok, err := o.IsValidExpression(expression)

	if ok {
		newExpr := o.ExpressionParser(expression)
		fmt.Println(newExpr.ListPriority, newExpr.ListSubExpr)
		vals.Add("expression", expression)
		resp, err := http.PostForm("http://localhost:9090/solvingExpression", vals)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Println(err)
	}

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