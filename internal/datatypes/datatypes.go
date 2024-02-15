package datatypes

type Status int // enum (перечесление) на подобие как в C#

const (
	Disable Status = iota
	Reconnect
	Idle
	Work
)

type Expression struct {
	ListPriority *map[int]int     `json:"listpriority"`
	ListSubExpr  *[]SubExpression `json:"listsubexpr"`
}

type SubExpression struct {
	Left     string `json:"left"`
	Right    string `json:"right"`
	Operator string `json:"operator"`
}

type Server struct {
	Url            string
	Status         Status
	CurrentTask    []int
	CountFailPings int
}

type Data struct {
	List     []any
	Settings map[string]int
	Status   Status
	Done     bool
}

type Task struct {
	Id              string
	Expression      SubExpression
	TimeExec        int
	IndexExpression int
	MaxIndex        int
}

func NewExpression(listPriority *map[int]int, listSubExpr *[]SubExpression) *Expression {
	return &Expression{
		ListPriority: listPriority,
		ListSubExpr:  listSubExpr,
	}
}
