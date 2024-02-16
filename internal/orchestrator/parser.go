package orchestrator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
)

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
				newSubExpr.Status = datatypes.Idle
				SubExpressions = append(SubExpressions, *newSubExpr)
				newSubExpr = &datatypes.SubExpression{Left: temp, Operator: ch}

				switch ch {
				case "*":
					newSubExpr.NameTimeExec = datatypes.TimeMulti
				case "/":
					newSubExpr.NameTimeExec = datatypes.TimeDivision
				case "+":
					newSubExpr.NameTimeExec = datatypes.TimeSum
				case "-":
					newSubExpr.NameTimeExec = datatypes.TimeSubtraction
				}

				temp = ""
				continue
			} else { // первое выражение
				if i == 0 {
					temp += ch
					continue
				}
				newSubExpr.Left = temp
				newSubExpr.Operator = ch

				switch ch {
				case "*":
					newSubExpr.NameTimeExec = datatypes.TimeMulti
				case "/":
					newSubExpr.NameTimeExec = datatypes.TimeDivision
				case "+":
					newSubExpr.NameTimeExec = datatypes.TimeSum
				case "-":
					newSubExpr.NameTimeExec = datatypes.TimeSubtraction
				}
				
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
	return *datatypes.NewExpression(ans, SubExpressions)
}

// сортировка по приоритету
func SortExpressions(SubExpressions []datatypes.SubExpression) map[int]datatypes.Node {
	answer := make(map[int]datatypes.Node)

	len := len(SubExpressions)

	priority := 0

	for i := 0; i < len; i++ { // сортировка для * и /
		if SubExpressions[i].Operator == "*" || SubExpressions[i].Operator == "/" {
			answer[priority] = datatypes.Node{Index: i, Status: datatypes.Idle}
			priority++
		}
	}

	for i := 0; i < len; i++ { // сортировка для + и -
		if SubExpressions[i].Operator == "+" || SubExpressions[i].Operator == "-" {
			answer[priority] = datatypes.Node{Index: i, Status: datatypes.Idle}
			priority++
		}
	}

	return answer
}
