package parser

type QueryPlan struct {
	Select SelectAction
	From   FromAction
	Where  *Filter
}

type SelectAction struct {
	Columns []string
}

type FromAction struct {
	Tables []string
}

type LogicalOperator uint8

const (
	LogicalOperatorAnd LogicalOperator = 1
	LogicalOperatorOr  LogicalOperator = 2
)

type Filter struct {
	Left            *Filter
	LogicalOperator LogicalOperator
	Right           *Filter
	Comparison      *Comparison
}

type Comparison struct {
	Left     string
	Operator string
	Right    string
}
