package orchestrator

import (
	"net/http"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
)

func (o *Orchestrator) CheckTasks() {
	for id, data := range o.ListExpr {
		for i, val := range data.ListPriority {
			if val.Agent != "" {
				_, err := http.Get(val.Agent)
				if err != nil {
					o.CancelTask(val.Agent, id, i)
				}
			}
		}
	}
}

func (o *Orchestrator) CancelTask(agentURL string, idExpr string, indexExpr int) {
	if idExpr == "" && indexExpr < 0 {
		for id, data := range o.ListExpr {
			for i, val := range data.ListPriority {
				if val.Agent == agentURL {
					val.Agent = ""
					val.Status = datatypes.Idle
					o.ListExpr[id].ListPriority[i] = val
					o.ListExpr[id].ListSubExpr[val.Index].Status = datatypes.Idle
					return
				}
			}
		}
	} else {
		val := o.ListExpr[idExpr].ListPriority[indexExpr]
		val.Agent = ""
		val.Status = datatypes.Idle
		o.ListExpr[idExpr].ListPriority[indexExpr] = val
		o.ListExpr[idExpr].ListSubExpr[indexExpr].Status = datatypes.Idle
	}
}

func (o *Orchestrator) GetTask(agentURL string) (datatypes.Task, bool) {
	for id, data := range o.ListExpr {
		for i := 0; i < len(data.ListPriority); i++ {
			expr := data.ListPriority[i]

			if expr.Status == datatypes.Idle {
				var newTask datatypes.Task
				otherUses := make([]int, 0, 2)
				if expr.Index == 0 { // если это первое выражение
					copyExpr := data.ListSubExpr[expr.Index]

					if expr.Index != len(data.ListSubExpr)-1 { // выражение правее
						if data.ListSubExpr[expr.Index+1].Status == datatypes.Work {
							continue
						}
						if expr.Index+1 < len(data.ListSubExpr) {
							if IsMultiOrDivision(data.ListSubExpr[expr.Index+1].Operator) && !IsMultiOrDivision(copyExpr.Operator) {
								if data.ListSubExpr[expr.Index+1].Answer == "" {
									continue
								}
								copyExpr.Right = data.ListSubExpr[expr.Index+1].Answer
								otherUses = append(otherUses, expr.Index+1)
							} else if data.ListSubExpr[expr.Index+1].Answer != "" {
								copyExpr.Right = data.ListSubExpr[expr.Index+1].Answer
								otherUses = append(otherUses, expr.Index+1)
							}
						}
					}

					data.ListSubExpr[expr.Index].Status = datatypes.Work
					copyExpr.Status = datatypes.Work
					newTask = *datatypes.NewTask(id, copyExpr, o.Settings[copyExpr.NameTimeExec], expr.Index, otherUses)
					expr.Status = datatypes.Work
					expr.Agent = agentURL
					data.ListPriority[i] = expr
					return newTask, true
				}

				copyExpr := data.ListSubExpr[expr.Index]

				// если левое или правое выражение обрабатывается, то переходим к другому
				if data.ListSubExpr[expr.Index-1].Status == datatypes.Work {
					if !IsMultiOrDivision(data.ListSubExpr[expr.Index-1].Operator) || !IsMultiOrDivision(copyExpr.Operator) {
						continue
					}
				}
				if expr.Index != len(data.ListSubExpr)-1 {
					if expr.Index+1 >= len(data.ListSubExpr) || data.ListSubExpr[expr.Index+1].Status == datatypes.Work {
						continue
					}
				}

				if IsMultiOrDivision(data.ListSubExpr[expr.Index-1].Operator) { // выражение ливее
					if data.ListSubExpr[expr.Index-1].Answer == "" {
						continue
					}
					copyExpr.Left = data.ListSubExpr[expr.Index-1].Answer
					otherUses = append(otherUses, expr.Index-1)

				} else if data.ListSubExpr[expr.Index-1].Answer != "" {
					copyExpr.Left = data.ListSubExpr[expr.Index-1].Answer
					otherUses = append(otherUses, expr.Index-1)
				}

				if expr.Index+1 < len(data.ListSubExpr) { // выражение правее
					if IsMultiOrDivision(data.ListSubExpr[expr.Index+1].Operator) && !IsMultiOrDivision(copyExpr.Operator) {
						if data.ListSubExpr[expr.Index+1].Answer == "" {
							continue
						}
						copyExpr.Right = data.ListSubExpr[expr.Index+1].Answer
						otherUses = append(otherUses, expr.Index+1)

					} else if data.ListSubExpr[expr.Index+1].Answer != "" {
						copyExpr.Right = data.ListSubExpr[expr.Index+1].Answer
						otherUses = append(otherUses, expr.Index+1)
					}
				}

				data.ListSubExpr[expr.Index].Status = datatypes.Work
				copyExpr.Status = datatypes.Work
				newTask = *datatypes.NewTask(id, copyExpr, o.Settings[copyExpr.NameTimeExec], expr.Index, otherUses)
				expr.Status = datatypes.Work
				expr.Agent = agentURL
				data.ListPriority[i] = expr
				return newTask, true
			}
		}
	}
	return datatypes.Task{}, false
}