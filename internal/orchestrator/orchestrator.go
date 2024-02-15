package orchestrator

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Orchestrator struct {
	Router      *mux.Router
	Tmpl        *template.Template
	ListExpr    []datatypes.Expression
	ListServers []datatypes.Server
	Settings    datatypes.Settings
	Data        datatypes.Data
}

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{
		Router:      mux.NewRouter(),
		Tmpl:        template.Must(template.ParseGlob("../../templates/*.html")),
		ListExpr:    make([]datatypes.Expression, 0),
		ListServers: make([]datatypes.Server, 0, 3),
	}
	return o
}

func (o *Orchestrator) InitRoutes() {
	o.Router.HandleFunc("/", o.handleGetIndex).Methods(http.MethodGet)

	o.Router.HandleFunc("/calculator", o.handleGetCalculator).Methods(http.MethodGet)
	o.Router.HandleFunc("/calculator", o.handlePostCalculator).Methods(http.MethodPost)

	o.Router.HandleFunc("/settings", o.handleGetSettings).Methods(http.MethodGet)
	o.Router.HandleFunc("/settings", o.handlePostSettings).Methods(http.MethodPost)

	o.Router.HandleFunc("/results", o.handleGetResult).Methods(http.MethodGet)
	o.Router.HandleFunc("/results", o.handlePostResult).Methods(http.MethodPost)

	o.Router.HandleFunc("/addServer", o.handlePostAddServer).Methods(http.MethodPost)

	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
	o.StartPingAgents()
}

// Запускает бесконечную горутину и каждые 5 секунд проверяет агентов на работоспасобность (мониторинг)
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

func (o *Orchestrator) IsValidExpression(s string) (bool, error) {
	s = strings.ToLower(s)
	if len(s) <= 2 { // выражение должно хотя бы быть формата "2+2"
		return false, errors.New("невалидное выражение, выражение слишком маленькое")
	}
	if strings.ContainsAny(s, "№!@#$%^&()~`qwertyuiop[]\\asdfghjkl;'zxcvbnm,.?йцукенгшщзхъфывапролджэячсмитьбю.|\":_ё=") {
		return false, errors.New("невалидное выражение, выражение содержит недопустимые символы")
	}

	temp := ""
	for i, ch := range s {
		switch temp {
		case "*", "/", "+", "-":
			switch string(ch) {
			case "*", "/", "+":
				return false, fmt.Errorf("невалидное выражение, недопускается \"%s%s\"", temp, string(s[i]))
			case "-":
				if i+1 < len(s) {
					switch string(s[i+1]) {
					case "*", "/", "+", "-":
						return false, fmt.Errorf("невалидное выражение, недопускается \"%s%s%s\"", temp, string(s[i]), string(s[i+1]))
					}
				}
				temp = string(ch)
			default:
				temp = string(ch)
			}
		default:
			temp = string(ch)
		}
	}
	switch string(s[len(s)-1]) {
	case "*", "/", "+", "-":
		return false, fmt.Errorf("невалидное выражение, в конце выражения не может быть \"%s\"", string(s[len(s)-1]))
	}
	return true, nil

}

// разбивает строку выражения на datatypes.SubExpression и возвращает тип datatypes.Expression
func (o *Orchestrator) ExpressionParser(s string) datatypes.Expression {
	s = strings.ReplaceAll(s, " ", "")

	chars := strings.Split(s, "")
	countOperators := 0
	for _, ch := range chars {
		switch ch {
		case "*", "/", "+", "-":
			countOperators++
		}
	}

	SubExpressions := make([]datatypes.SubExpression, 0, countOperators)
	newSubExpr := &datatypes.SubExpression{}
	temp := ""
	notFirst := false

	for i, ch := range chars {
		if ch == "+" || ch == "-" || ch == "*" || ch == "/" {
			if notFirst {
				switch chars[i-1] {
				case "+", "-", "/", "*":
					temp += ch
					continue
				}
				newSubExpr.Right = temp
				SubExpressions = append(SubExpressions, *newSubExpr)
				newSubExpr = &datatypes.SubExpression{Left: temp, Operator: ch}
				temp = ""
				continue
			} else { // первое выражение
				if i == 0 {
					temp += ch
					continue
				}
				newSubExpr.Left = temp
				newSubExpr.Operator = ch
				temp = ""
				notFirst = true
				continue
			}
		}
		temp += ch
	}
	newSubExpr.Right = temp
	SubExpressions = append(SubExpressions, *newSubExpr)

	ans := SortExpressions(SubExpressions)
	return *datatypes.NewExpression(&ans, &SubExpressions)
}

// сортировка по приоритету
func SortExpressions(SubExpressions []datatypes.SubExpression) map[int]int {
	answer := make(map[int]int)

	len := len(SubExpressions)

	priority := 0

	for i := 0; i < len; i++ { // сортировка для * и /
		if SubExpressions[i].Operator == "*" || SubExpressions[i].Operator == "/" {
			answer[priority] = i
			priority++
		}
	}

	for i := 0; i < len; i++ { // сортировка для + и -
		if SubExpressions[i].Operator == "+" || SubExpressions[i].Operator == "-" {
			answer[priority] = i
			priority++
		}
	}

	return answer
}
