package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		if err = database.AddExpression(expression, &newExpr, "Work", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			fmt.Println(err)
			return
		}
		o.ListExpr[expression] = &newExpr
	} else {
		fmt.Println(err)
		if err = database.AddExpression(expression, nil, "Fail", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			fmt.Println(err)
			return
		}
	}

}

// Settings.html
func (o *Orchestrator) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	listSettings := make(map[string]int)
	for key, val := range o.Settings {
		listSettings[string(key)] = val
	}
	o.Data.Settings = listSettings
	o.Tmpl.ExecuteTemplate(w, "settings.html", o.Data)
	o.Data = datatypes.Data{}
}

func (o *Orchestrator) handlePostSettings(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/settings", http.StatusSeeOther)
	targetName := make([]datatypes.NameTimeExec, 0, 5)
	targetName = append(targetName, datatypes.TimeSum, datatypes.TimeSubtraction, datatypes.TimeMulti, datatypes.TimeDivision, datatypes.TimeOut)
	for _, valname := range targetName {
		val := r.PostFormValue(string(valname))
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
			o.Settings[valname] = num
		}
	}
	err := database.UpdateSettingsData(o.Settings)
	if err != nil {
		fmt.Println(err)
	}
	o.Data.Done = true
}

// Results.html
func (o *Orchestrator) handleGetResult(w http.ResponseWriter, r *http.Request) {
	list, err := database.GetAllExpression()
	if err != nil {
		fmt.Println(err)
		o.Tmpl.ExecuteTemplate(w, "results.html", nil)
	} else {
		o.Data.List = list
		o.Tmpl.ExecuteTemplate(w, "results.html", o.Data)
	}
}

func (o *Orchestrator) handlePostAddServer(w http.ResponseWriter, r *http.Request) {
	val := r.PostFormValue("server")
	for _, serv := range o.ListServers {
		if serv.Url == val {
			return
		}
	}
	o.ListServers = append(o.ListServers, datatypes.Server{Url: val, Status: datatypes.Idle})
}

func (o *Orchestrator) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	isAgent := false
	url := r.URL.Query().Get("agent")

	
	for _, serv := range o.ListServers {
		if url == serv.Url {
			isAgent = true
			break
		}
	}

	if isAgent {
		newTask, ok := o.GetTask(url)
		if !ok {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(newTask); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (o *Orchestrator) handlePostAnswer(w http.ResponseWriter, r *http.Request) {
	var task datatypes.Task
	ans := r.PostFormValue("answer")
	err := json.Unmarshal([]byte(ans), &task)
	if err != nil {
		fmt.Println(err)
		return
	}
	go o.CheckAndUpdateExpression(task)
}
