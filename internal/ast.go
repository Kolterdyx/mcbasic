package internal

type Expr interface {
}

type Binary struct {
	Expr

	Left     Expr
	Operator Token
	Right    Expr
}

type Grouping struct {
	Expr

	Expression Expr
}

type Literal struct {
	Expr

	Value interface{}
}
