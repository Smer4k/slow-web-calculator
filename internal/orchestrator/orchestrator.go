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

func (o *Orchestrator) IsValidExpression(s string) (bool, error) {
	s = strings.ToLower(s)
	if len(s) <= 2 { // выражение должно хотя бы быть формата "2+2"
		return false, errors.New("Невалидное выражение, выражение слишком маленькое")
	}
	if strings.ContainsAny(s, "№!@#$%^&()~`qwertyuiop[]\\asdfghjkl;'zxcvbnm,.?йцукенгшщзхъфывапролджэячсмитьбю.|\":_ё=") {
		return false, errors.New("Выражение содержит недопустимые символы")
	}
	temp := ""
	for _, ch := range s { // <----- здесь баг
		switch temp {
		case "*", "/", "+", "-": // переделать, мб нужно добавить возможность использовать "2 + -3"
			switch string(ch) {
			case "*", "/", "+", "-":
				return false, errors.New("Выражение неправильного формата, нельзя чтобы шло два и более \"++\" или \"+-\" и т.д")
			}
		default:
			temp = string(ch)
		}
	}
	return true, nil
}

// разбивает строку выражения на datatypes.SubExpression и возвращает тип datatypes.Expression
func (o *Orchestrator) ExpressionParser(s string) *datatypes.Expression {
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

	for _, ch := range chars {
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
	}
	newSubExpr.Right = temp
	SubExpressions = append(SubExpressions, *newSubExpr)

	ans := SortExpressions(SubExpressions)
	return datatypes.NewExpression(&ans, &SubExpressions)
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
