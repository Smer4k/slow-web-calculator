package datatypes

type Expression struct {
	List  *map[int]SubExpression
	Count uint
}

type SubExpression struct {
	Left string
	Right string
	Operator string
	Answer int
}

func NewExpression(list *map[int]SubExpression, count uint) *Expression {
	return &Expression{
		List:  list,
		Count: count,
	}
}
