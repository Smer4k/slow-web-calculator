package datatypes

import "time"

type Status int // enum (перечесление) на подобие как в C#
type NameTimeExec string

const (
	Disable Status = iota
	Reconnect
	Idle
	Work
	Done
	BadRequest
	ServerError
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
	LastPing       time.Time
	CancelDelChan  chan struct{}
}

type Data struct {
	Settings map[string]int
	Status   Status
	Done     string
	Text     string
}

type DataServer struct {
	Status   string
	TimePing string
}

type DataExpression struct {
	Expr      string
	Answer    string
	Status    string
	TimeSend  string
	TimeSolve string
}

type Task struct {
	Id              string        `json:"id"`
	Expression      SubExpression `json:"expression"`
	TimeExec        int           `json:"timeexec"`
	IndexExpression int           `json:"indexexpression"`
	Answer          string       `json:"answer"`
	OtherUses       []int         `json:"otheruses"`
}

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

func NewTask(expressionOrigin string, expression SubExpression, timeExec, index int, otherUses []int) *Task {
	return &Task{
		Id:              expressionOrigin,
		Expression:      expression,
		TimeExec:        timeExec,
		IndexExpression: index,
		OtherUses:       otherUses,
	}
}
