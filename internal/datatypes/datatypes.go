package datatypes

type Status int // enum (перечесление) на подобие как в C#
type NameTimeExec string

const (
	Disable Status = iota
	Reconnect
	Idle
	Work
	Done
)

const (
	TimeSum         NameTimeExec = "time_sum"
	TimeSubtraction NameTimeExec = "time_subtraction"
	TimeMulti       NameTimeExec = "time_multi"
	TimeDivision    NameTimeExec = "time_division"
	TimeOut         NameTimeExec = "time_out"
)

type Expression struct {
	ListPriority map[int]Node    `json:"listpriority"`
	ListSubExpr  []SubExpression `json:"listsubexpr"`
}

type SubExpression struct {
	Left         string       `json:"left"`
	Right        string       `json:"right"`
	Operator     string       `json:"operator"`
	Answer       string       `json:"answer"`
	NameTimeExec NameTimeExec `json:"nametimeexec"`
	Status       Status       `json:"status"`
}

type Server struct {
	Url            string
	Status         Status
	CountFailPings int
}

type Data struct {
	List     []any
	Settings map[NameTimeExec]int
	Status   Status
	Done     bool
}

type Task struct {
	Id              string        `json:"id"`
	Expression      SubExpression `json:"expression"`
	TimeExec        int           `json:"timeexec"`
	IndexExpression int           `json:"indexexpression"`
}

// у меня уже закончились идее как назвать
type Node struct {
	Index  int    `json:"index"`
	Agent  string `json:"agent"`
	Status Status `json:"status"`
}

func NewExpression(listPriority map[int]Node, listSubExpr []SubExpression) *Expression {
	return &Expression{
		ListPriority: listPriority,
		ListSubExpr:  listSubExpr,
	}
}

func NewTask(expressionOrigin string, expression SubExpression, timeExec, index int) *Task {
	return &Task{
		Id:              expressionOrigin,
		Expression:      expression,
		TimeExec:        timeExec,
		IndexExpression: index,
	}
}
