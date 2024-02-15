package orchestrator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Smer4k/slow-web-calculator/internal/database"
	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
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
	expression := r.FormValue("expression")
	ok, err := o.IsValidExpression(expression)

	if ok {
		newExpr := o.ExpressionParser(expression)
		if err = database.AddExpression(expression, &newExpr, 2, "work"); err != nil {
			fmt.Println(err)
			return
		}
		o.ListExpr = append(o.ListExpr, newExpr)
	} else {
		fmt.Println(err)
		if err = database.AddExpression(expression, nil, 0, "fail"); err != nil {
			fmt.Println(err)
			return
		}
	}

}

// Settings.html
func (o *Orchestrator) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	o.Data.Settings = o.Settings
	o.Tmpl.ExecuteTemplate(w, "settings.html", o.Data)
	o.Data = datatypes.Data{}
}

func (o *Orchestrator) handlePostSettings(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/settings", http.StatusSeeOther)
	targetName := make([]string, 0, 5)
	targetName = append(targetName, "time_+", "time_-", "time_*", "time_/", "time_out")
	for _, valname := range targetName {
		val := r.PostFormValue(valname)
		if val != "" {
			num, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if num < 0 {
				fmt.Printf("\"%s\" не может быть меньше 0", valname)
				continue
			}
			switch valname {
			case "time_+":
				o.Settings.TimeSum = num
			case "time_-":
				o.Settings.TimeDeduction = num
			case "time_*":
				o.Settings.TimeMulti = num
			case "time_/":
				o.Settings.TimeDivision = num
			case "time_out":
				o.Settings.TimeDisplayServer = num
			}
		}
	}
	o.Data.Done = true
}

// Results.html
func (o *Orchestrator) handleGetResult(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "results.html", o.Data)
}

func (o *Orchestrator) handlePostResult(w http.ResponseWriter, r *http.Request) {

}

func (o *Orchestrator) handlePostAddServer(w http.ResponseWriter, r *http.Request) {
	val := r.PostFormValue("server")
	for _, serv := range o.ListServers {
		if serv.Url == val {
			return
		}
	}
	o.ListServers = append(o.ListServers, datatypes.Server{Url: val, Status: 1})
}
