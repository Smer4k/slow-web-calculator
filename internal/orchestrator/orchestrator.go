package orchestrator

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/database"
	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Orchestrator struct {
	Router      *mux.Router
	Tmpl        *template.Template
	ListExpr    map[string]datatypes.Expression // Список всех задач для агентов
	ListServers []datatypes.Server // Список подключенных агентов
	Settings    map[string]int
	Data        datatypes.Data
}

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{
		Router:      mux.NewRouter(),
		Tmpl:        template.Must(template.ParseGlob("../../templates/*.html")),
		ListExpr:    make(map[string]datatypes.Expression),
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
	o.Router.HandleFunc("/results", o.handlePostResult).Methods(http.MethodPost)

	o.Router.HandleFunc("/addServer", o.handlePostAddServer).Methods(http.MethodPost)
	o.Router.HandleFunc("/ping", o.handlePing)
	
	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
	o.StartPingAgents()
}

// Загружает данные с базы данных
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

// Запускает бесконечную горутину и каждые 10 секунд проверяет агентов на работоспасобность (мониторинг)
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
						o.ListServers = append(o.ListServers[:i], o.ListServers[i+1:]...)
						i--
					}
				} else {
					o.ListServers[i].CountFailPings = 0
					fmt.Printf("Сервер %s работает исправно.\n", o.ListServers[i].Url)
				}
			}
		}
	}()
}

