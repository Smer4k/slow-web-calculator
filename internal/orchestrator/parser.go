package orchestrator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
)

func (o *Orchestrator) IsValidExpression(s string) (bool, error) {
	if len(s) <= 2 { // выражение должно хотя бы быть формата "2+2"
		return false, errors.New("невалидное выражение, выражение слишком маленькое")
	}
	if strings.ContainsAny(s, "№!@#$%^&()~`qwertyuiop[]\\asdfghjkl;'zxcvbnm,?йцукенгшщзхъфывапролджэячсмитьбю|\":_ё=") {
		return false, errors.New("невалидное выражение, выражение содержит недопустимые символы")
	}

	temp := ""
	for i, ch := range s {
		switch temp {
		case "*", "/", "+", "-":
			switch string(ch) {
			case "*", "/", "+": // если два идет два подрят ** или подобие (искл. '-')
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
	switch string(s[len(s)-1]) { // если в конце окажется не цифра
	case "*", "/", "+", "-":
		return false, fmt.Errorf("невалидное выражение, в конце выражения не может быть \"%s\"", string(s[len(s)-1]))
	}
	return true, nil

}

func (o *Orchestrator) ExpressionParser(s string) datatypes.Expression {
	s = strings.ReplaceAll(s, " ", "")

	chars := strings.Split(s, "")
	countOperators := 0
	for _, ch := range chars { // смотрим сколько нужно места для массива (есть погрешность)
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
				switch chars[i-1] { // если этот ch является '-' и до него ch является оператором
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
				if i == 0 { // если первый ch это оператор
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
	newSubExpr.Status = datatypes.Idle
	SubExpressions = append(SubExpressions, *newSubExpr)

	ans := SortExpressions(SubExpressions)
	return *datatypes.NewExpression(ans, SubExpressions)
}

func SortExpressions(SubExpressions []datatypes.SubExpression) map[int]datatypes.Node {
	answer := make(map[int]datatypes.Node)

	len := len(SubExpressions)

	priority := 0

	for i := 0; i < len; i++ { // сортировка для * и /
		if IsMultiOrDivision(SubExpressions[i].Operator) {
			answer[priority] = datatypes.Node{Index: i, Status: datatypes.Idle}
			priority++
		}
	}

	for i := 0; i < len; i++ { // сортировка для + и -
		if !IsMultiOrDivision(SubExpressions[i].Operator) {
			answer[priority] = datatypes.Node{Index: i, Status: datatypes.Idle}
			priority++
		}
	}

	return answer
}
