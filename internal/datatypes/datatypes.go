package datatypes

type Expression struct {
	ListPriority *map[int]int `json:"listpriority"`
	ListSubExpr  *[]SubExpression `json:"listsubexpr"`
}

type SubExpression struct {
	Left     string `json:"left"`
	Right    string	`json:"right"`
	Operator string	`json:"operator"`
}

type Server struct {
	Url    string
	Status int
}

func NewExpression(listPriority *map[int]int, listSubExpr *[]SubExpression) *Expression {
	return &Expression{
		ListPriority: listPriority,
		ListSubExpr:  listSubExpr,
	}
}
