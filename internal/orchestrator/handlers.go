package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/database"
	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
)

func (o *Orchestrator) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	o.Tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Calculator.html
func (o *Orchestrator) handleGetCalculator(w http.ResponseWriter, r *http.Request) {
	if o.Data.Status == datatypes.BadRequest {
		w.WriteHeader(http.StatusBadRequest)
	} else if o.Data.Status == datatypes.ServerError {
		if strings.Contains(o.Data.Text, "UNIQUE") {
			o.Data.Text = "Выражение уже было ранее отправильно"
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	o.Tmpl.ExecuteTemplate(w, "calculator.html", o.Data)
	o.Data = datatypes.Data{}
}

func (o *Orchestrator) handlePostCalculator(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/calculator", http.StatusSeeOther)
	expression := r.FormValue("expression")
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, "×", "*")
	expression = strings.ReplaceAll(expression, "÷", "/")
	ok, err := o.IsValidExpression(expression)

	if ok {
		newExpr := o.ExpressionParser(expression)
		if err = database.AddExpression(expression, &newExpr, "Work", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			fmt.Println(err)
			o.Data.Text = err.Error()
			o.Data.Done = "danger"
			o.Data.Status = datatypes.ServerError
			return
		}
		o.Data.Done = "sucsess"
		o.ListExpr[expression] = &newExpr
	} else {
		fmt.Println(err)
		o.Data.Text = err.Error()
		o.Data.Done = "danger"
		o.Data.Status = datatypes.BadRequest
		if err = database.AddExpression(expression, nil, "Fail", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			fmt.Println(err)
			o.Data.Text = err.Error()
			o.Data.Done = "danger"
			o.Data.Status = datatypes.ServerError
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
	o.Data.Done = "sucsess"
}

// Results.html
func (o *Orchestrator) handleGetResult(w http.ResponseWriter, r *http.Request) {
	list, err := database.GetAllExpression()
	slices.SortFunc(list, func(a *datatypes.DataExpression, b *datatypes.DataExpression) int {
		timeA, _ := time.Parse("2006-01-02 15:04:05", a.TimeSend)
		timeB, _ := time.Parse("2006-01-02 15:04:05", b.TimeSend)
		if timeA.Before(timeB) {
			return -1
		} else if timeA.Equal(timeB) {
			return 0
		} else {
			return 1
		}
	})
	if err != nil {
		fmt.Println(err)
		o.Tmpl.ExecuteTemplate(w, "results.html", nil)
	} else {
		o.Tmpl.ExecuteTemplate(w, "results.html", list)
	}
}

func (o *Orchestrator) handleGetResources(w http.ResponseWriter, r *http.Request) {
	if len(o.ListServers) == 0 {
		o.Tmpl.ExecuteTemplate(w, "resources.html", nil)
		return
	}
	list := make(map[string]datatypes.DataServer)
	for _, serv := range o.ListServers {
		data := &datatypes.DataServer{}
		switch serv.Status {
		case datatypes.Idle:
			data.Status = "success"
		case datatypes.Reconnect:
			data.Status = "info"
		case datatypes.Disable:
			data.Status = "danger"
		}
		data.TimePing = serv.LastPing.Format("2006-01-02 15:04:05")
		list[serv.Url] = *data
	}
	o.Tmpl.ExecuteTemplate(w, "resources.html", list)
}

// Ниже обработчики для агентов

func (o *Orchestrator) handlePostAddServer(w http.ResponseWriter, r *http.Request) {
	val := r.PostFormValue("server")
	for i, serv := range o.ListServers {
		if serv.Url == val {
			if serv.Status == datatypes.Disable || serv.Status == datatypes.Reconnect {
				o.ListServers[i].CountFailPings = 0
				o.ListServers[i].Status = datatypes.Idle
				o.ListServers[i].LastPing = time.Now()
				if serv.Status == datatypes.Disable {
					close(o.ListServers[i].CancelDelChan)
				}
			}
			return
		}
	}
	o.ListServers = append(o.ListServers, datatypes.Server{Url: val, Status: datatypes.Idle, CountFailPings: 0, LastPing: time.Now()})
	go o.StartPingAgent(o.ListServers[len(o.ListServers)-1].Url)
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
