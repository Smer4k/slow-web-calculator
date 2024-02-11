package datatypes

type Expression struct {
	ListPriority  *map[int]int
	ListSubExpr *[]SubExpression
}

type SubExpression struct {
	Left string
	Right string
	Operator string
}

func NewExpression(listPriority *map[int]int, listSubExpr *[]SubExpression) *Expression {
	return &Expression{
		ListPriority:  listPriority,
		ListSubExpr: listSubExpr,
	}
}
