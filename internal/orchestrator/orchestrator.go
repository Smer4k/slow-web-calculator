package orchestrator

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/database"
	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Orchestrator struct {
	Router      *mux.Router
	Tmpl        *template.Template
	ListExpr    map[string]*datatypes.Expression // Список всех задач для агентов
	ListServers []datatypes.Server               // Список подключенных агентов
	Settings    map[datatypes.NameTimeExec]int   // настройки сервера
	Data        datatypes.Data                   // отправляемые данные при запросе пользователя
}

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{
		Router:      mux.NewRouter(),
		Tmpl:        template.Must(template.ParseGlob("../../templates/*.html")),
		ListExpr:    make(map[string]*datatypes.Expression),
		ListServers: make([]datatypes.Server, 0, 3),
	}
	return o
}

func (o *Orchestrator) InitRoutes() {
	o.LoadData()
	o.Router.HandleFunc("/", o.handleGetIndex).Methods(http.MethodGet)

	o.Router.HandleFunc("/calculator", o.handleGetCalculator).Methods(http.MethodGet)
	o.Router.HandleFunc("/calculator", o.handlePostCalculator).Methods(http.MethodPost)

	o.Router.HandleFunc("/settings", o.handleGetSettings).Methods(http.MethodGet)
	o.Router.HandleFunc("/settings", o.handlePostSettings).Methods(http.MethodPost)

	o.Router.HandleFunc("/results", o.handleGetResult).Methods(http.MethodGet)

	o.Router.HandleFunc("/addServer", o.handlePostAddServer).Methods(http.MethodPost)
	o.Router.HandleFunc("/getExpression", o.handleGetExpression)
	o.Router.HandleFunc("/postAnswer", o.handlePostAnswer).Methods(http.MethodPost)

	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
	o.StartPingAgents()
}

func IsMultiOrDivision(operator string) bool {
	return strings.ContainsAny(operator, "*/")
}

func (o *Orchestrator) SetStatusForNeighbors(id string,index int) {
	for i := index; i >= 0; i-- {
		if !IsMultiOrDivision(o.ListExpr[id].ListSubExpr[i].Operator) {
			break
		}
		o.ListExpr[id].ListSubExpr[i].Status = datatypes.Done
	}
}

func (o *Orchestrator) LoadData() {
	settings, err := database.GetSettingsData()
	if err != nil {
		panic(err)
	}
	o.Settings = settings
	listExpr, err := database.GetWorkExpressionsData()
	if err != nil {
		panic(err)
	}
	o.ListExpr = listExpr
}

func (o *Orchestrator) StartPingAgents() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if len(o.ListServers) == 0 {
				fmt.Println("Нет подключенных агентов.")
				continue
			}

			for i := 0; i < len(o.ListServers); i++ {
				_, err := http.Get(o.ListServers[i].Url)

				if err != nil {
					o.ListServers[i].CountFailPings++
					fmt.Println(err)

					if o.ListServers[i].CountFailPings >= 3 {
						fmt.Printf("Сервер %s слишком долго не отвечал и был удален.\n", o.ListServers[i].Url)
						o.CancelTask(o.ListServers[i].Url, "", -1)
						o.ListServers = append(o.ListServers[:i], o.ListServers[i+1:]...)
						i--
					}
				} else if o.ListServers[i].CountFailPings != 0 {
					o.ListServers[i].CountFailPings = 0
					fmt.Printf("Сервер %s работает исправно.\n", o.ListServers[i].Url)
				}
			}
		}
	}()
}

func (o *Orchestrator) CheckAndUpdateExpression(task datatypes.Task) {
	for {
		if len(o.ListExpr) != 0 {
			break
		}
	}
	o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Answer = strconv.FormatFloat(task.Answer, 'f', -1, 64)
	for key, val := range o.ListExpr[task.Id].ListPriority {
		if val.Index == task.IndexExpression {
			if task.IndexExpression == len(o.ListExpr[task.Id].ListSubExpr)-1 {
				val.Status = datatypes.Done
				o.SetStatusForNeighbors(task.Id, task.IndexExpression-1)

			} else if IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Operator) && IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression+1].Operator) {
				val.Status = datatypes.Work
				o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Status = datatypes.Work

			} else {
				val.Status = datatypes.Done
				o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Status = datatypes.Done
			}

			val.Agent = ""
			o.ListExpr[task.Id].ListPriority[key] = val
			break
		}
	}

	if IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Operator) {
		for i := task.IndexExpression - 1; i >= 0; i-- {
			val := o.ListExpr[task.Id].ListSubExpr[i]
			if IsMultiOrDivision(val.Operator) {
				if val.Answer != "" {
					o.ListExpr[task.Id].ListSubExpr[i].Answer = strconv.FormatFloat(task.Answer, 'f', -1, 64)
				} else {
					break
				}
			} else {
				break
			}
		}
	} else {
		for _, val := range task.OtherUses {
			o.ListExpr[task.Id].ListSubExpr[val].Answer = strconv.FormatFloat(task.Answer, 'f', -1, 64)
		}
	}

	lastIndex := o.ListExpr[task.Id].ListPriority[len(o.ListExpr[task.Id].ListPriority)-1].Index
	if o.ListExpr[task.Id].ListSubExpr[lastIndex].Answer != "" {
		err := database.UpdateExpression(task.Id, nil, "Done", o.ListExpr[task.Id].ListSubExpr[lastIndex].Answer, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println(err)
			return
		}
		delete(o.ListExpr, task.Id)
	} else {
		database.UpdateExpression(task.Id, o.ListExpr[task.Id], "Work", "", "")
	}
}

