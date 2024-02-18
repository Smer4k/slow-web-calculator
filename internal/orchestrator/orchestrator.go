package orchestrator

import (
	"fmt"
	"html/template"
	"math/rand"
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
	o.Router.HandleFunc("/resources", o.handleGetResources)

	o.Router.HandleFunc("/addServer", o.handlePostAddServer).Methods(http.MethodPost)
	o.Router.HandleFunc("/getExpression", o.handleGetExpression)
	o.Router.HandleFunc("/postAnswer", o.handlePostAnswer).Methods(http.MethodPost)

	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
}

func IsMultiOrDivision(operator string) bool {
	return strings.ContainsAny(operator, "*/")
}

func (o *Orchestrator) SetStatusNeighborsMultiDivision(id string, index int) {
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
	for key, val := range o.ListExpr {
		fmt.Println(key, " ", *val)
	}
}

func (o *Orchestrator) StartPingAgent(agentURL string) {
	seconds := 10 * time.Second
	randSecond := time.Duration(float64(time.Second) * rand.Float64())
	ticker := time.NewTicker(seconds + randSecond)
	defer ticker.Stop()
	for range ticker.C {
		for i := 0; i < len(o.ListServers); i++ {
			if o.ListServers[i].Url != agentURL {
				continue
			}
			if o.ListServers[i].Status == datatypes.Disable {
				continue
			}
			_, err := http.Get(agentURL)
			o.ListServers[i].LastPing = time.Now()
			if err != nil {
				o.ListServers[i].CountFailPings++
				if o.ListServers[i].Status != datatypes.Reconnect {
					o.ListServers[i].Status = datatypes.Reconnect
					o.CancelTask(agentURL, "", -1)
					fmt.Printf("Сервер %s не отвечает, все задачи сервера сняты, ошибка:\n%s\n", o.ListServers[i].Url, err)
				}

				if o.ListServers[i].CountFailPings >= 3 {
					o.ListServers[i].Status = datatypes.Disable
					fmt.Printf("Сервер %s слишком долго не отвечал и был поставлен на удаление\n", agentURL)
					
					o.ListServers[i].CancelDelChan = make(chan struct{})
					go o.DeleteServer(agentURL, o.ListServers[i].CancelDelChan)
				}
			} else if o.ListServers[i].CountFailPings != 0 {
				o.ListServers[i].CountFailPings = 0
				fmt.Printf("Сервер %s работает исправно.\n", agentURL)
			}
		}
	}
}

func (o *Orchestrator) CheckAndUpdateExpression(task datatypes.Task) {
	for { // ждем пока сервер не подгузит выражения из БД
		if len(o.ListExpr) != 0 {
			break
		}
	}

	o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Answer = strconv.FormatFloat(task.Answer, 'f', -1, 64)
	
	for key, val := range o.ListExpr[task.Id].ListPriority {
		if val.Index == task.IndexExpression {
			if task.IndexExpression == len(o.ListExpr[task.Id].ListSubExpr)-1 { // если это последние выражение
				val.Status = datatypes.Done
				o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Status = datatypes.Done
				o.SetStatusNeighborsMultiDivision(task.Id, task.IndexExpression-1)

			} else if IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Operator) && IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression+1].Operator) {
				// если это и левое выражение является * или /
				val.Status = datatypes.Work
				o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Status = datatypes.Work
			} else if len(task.OtherUses) != 0 && !IsMultiOrDivision(o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Operator) {
				val.Status = datatypes.Done
				o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Status = datatypes.Done
				for _, val := range task.OtherUses {
					o.ListExpr[task.Id].ListSubExpr[val].Status = datatypes.Done
				}
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
			if IsMultiOrDivision(val.Operator) { // заменяет ответ у соседних * и / на новый ответ
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
		for _, val := range task.OtherUses { // заменяет ответ у левого и правого выражения
			o.ListExpr[task.Id].ListSubExpr[val].Answer = strconv.FormatFloat(task.Answer, 'f', -1, 64)
		}
	}

	allSolve := true
	for _, val := range o.ListExpr[task.Id].ListSubExpr {
		if val.Answer == "" {
			allSolve = false
			break
		}
	}
	if allSolve { // если последние выражение было решено (по списку приоритета)
		err := database.UpdateExpression(task.Id, nil, "Done", o.ListExpr[task.Id].ListSubExpr[task.IndexExpression].Answer, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println(err)
			return
		}
		delete(o.ListExpr, task.Id)
	} else {
		database.UpdateExpression(task.Id, o.ListExpr[task.Id], "Work", "", "")
	}
}

// Удаляет сервер из списка через время
func (o *Orchestrator) DeleteServer(agentURL string, cancel chan struct{}) {
	seconds := time.Duration(o.Settings[datatypes.TimeOut]) * time.Second
	randSecond := time.Duration(float64(time.Second) * rand.Float64())
	time := time.After(seconds + randSecond)
	for {
		select {
		case <-time:
			for i := 0; i < len(o.ListServers); i++ {
				if o.ListServers[i].Url != agentURL {
					continue
				}
				o.ListServers = append(o.ListServers[:i], o.ListServers[i+1:]...)
			}
			return
		case <-cancel:
			return
		}
	}
}
