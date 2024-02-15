package datatypes

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
	Status         int
	CountFailPings int
}

type Settings struct {
	TimeSum           int
	TimeDeduction     int
	TimeMulti         int
	TimeDivision      int
	TimeDisplayServer int
}

type Data struct {
	List     []any
	Settings Settings
	Status   string
	Done     bool
}

func NewExpression(listPriority *map[int]int, listSubExpr *[]SubExpression) *Expression {
	return &Expression{
		ListPriority: listPriority,
		ListSubExpr:  listSubExpr,
	}
}
