package parser

type QueryPlan struct {
	Select SelectAction `json:"select"`
	From   FromAction   `json:"from"`
	Where  *Filter      `json:"where,omitempty"`
}

type SelectAction struct {
	Columns []string `json:"columns"`
}

type FromAction struct {
	Tables []string `json:"tables"`
}

type LogicalOperator string

const (
	LogicalOperatorAnd LogicalOperator = "AND"
	LogicalOperatorOr  LogicalOperator = "OR"
)

type Filter struct {
	Left            *Filter         `json:"left,omitempty"`
	LogicalOperator LogicalOperator `json:"logicalOperator,omitempty"`
	Right           *Filter         `json:"right,omitempty"`
	Comparison      *Comparison     `json:"comparison,omitempty"`
}

type Comparison struct {
	Left     string `json:"left"`
	Operator string `json:"operator"`
	Right    string `json:"right"`
}
