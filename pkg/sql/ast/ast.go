package ast

type Statement interface {
	statement()
}

type SelectStatement struct {
	Columns []Expr
	From    string
	Where   Expr // nil when WHERE is absent
}

func (*SelectStatement) statement() {}

type Expr interface {
	expr()
}

type ColumnRef struct{ Name string }

func (*ColumnRef) expr() {}

type IntLiteral struct{ Value int64 }

func (*IntLiteral) expr() {}

type Star struct{}

func (*Star) expr() {}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (*BinaryExpr) expr() {}

type StringLiteral struct{ Value string }

func (*StringLiteral) expr() {}

type InsertStatement struct {
	Table   string
	Columns []string // nil when omitted
	Values  [][]Expr // one entry per VALUES row
}

func (*InsertStatement) statement() {}

type Assignment struct {
	Column string
	Value  Expr
}

type UpdateStatement struct {
	Table       string
	Assignments []Assignment
	Where       Expr
}

func (*UpdateStatement) statement() {}

type DeleteStatement struct {
	Table string
	Where Expr
}

func (*DeleteStatement) statement() {}

type CreateDatabaseStatement struct{ Name string }

func (*CreateDatabaseStatement) statement() {}

type DropDatabaseStatement struct{ Name string }

func (*DropDatabaseStatement) statement() {}
