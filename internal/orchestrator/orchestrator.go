package orchestrator

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Orchestrator struct {
	Router *mux.Router
	Tmpl   *template.Template
	Data   any
}

func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{
		Router: mux.NewRouter(),
		Tmpl:   template.Must(template.ParseGlob("../../templates/*.html")),
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
	o.Router.HandleFunc("/postResult", o.handlePostResult).Methods(http.MethodPost)
	http.Handle("/", o.Router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../templates/static/"))))
}

func (o *Orchestrator) ExpressionParser(s string) (*datatypes.Expression, error) {
	s = strings.ReplaceAll(s, " ", "")
	lenS := len(s)
	if lenS <= 2 { // выражение должно хотя бы быть формата "2+2"
		return nil, errors.New("Невалидное выражение, выражение слишком маленькое")
	}

	chars := strings.Split(s, "")
	SubExpressions := make([]datatypes.SubExpression, 0, 100)
	newSubExpr := &datatypes.SubExpression{}
	temp := ""
	notFirst := false

	for i, ch := range chars {
		if ch == "+" || ch == "-" || ch == "*" || ch == "/" {
			if notFirst {
				newSubExpr.Right = temp
				SubExpressions = append(SubExpressions, *newSubExpr)
				newSubExpr = &datatypes.SubExpression{Left: temp, Operator: ch}
				temp = ""
				continue
			} else { // первое выражение
				newSubExpr.Left = temp
				newSubExpr.Operator = ch
				temp = ""
				notFirst = true
				continue
			}
		}
		temp += ch
		if i == lenS-1 { // добавляем последний элемент
			newSubExpr.Right = temp
			SubExpressions = append(SubExpressions, *newSubExpr)
		}
	}
	ans := SortExpression(SubExpressions)
	return datatypes.NewExpression(&ans, uint(len(ans))), nil
}

// переделать
func SortExpression(SubExpressions []datatypes.SubExpression) map[int]datatypes.SubExpression {
	answer := make(map[int]datatypes.SubExpression)

	len := len(SubExpressions)
	priority := 0

	for i := 0; i < len; i++ {
		switch SubExpressions[i].Operator {
		case "*":
			if i == len-1 {
				SubExpressions[i-1].Right = "?"
				answer[priority] = SubExpressions[i]
				priority++
				answer[priority] = SubExpressions[i-1]
				priority++
			} else if i == 0 {
				SubExpressions[i+1].Left = "?"
				answer[priority] = SubExpressions[i]
				priority++
				answer[priority] = SubExpressions[i+1]
				priority++
			}
		case "/":
			if i == len-1 {
				SubExpressions[i-1].Right = "?"
				answer[priority] = SubExpressions[i]
				priority++
				answer[priority] = SubExpressions[i-1]
				priority++
			} else if i == 0 {
				SubExpressions[i+1].Left = "?"
				answer[priority] = SubExpressions[i]
				priority++
				answer[priority] = SubExpressions[i+1]
				priority++
			} else {
				SubExpressions[i+1].Left = "?"
				SubExpressions[i-1].Right = "?"
				answer[priority] = SubExpressions[i]
				priority++
			}
		}
	}
	return answer
}