package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bengosborn/stream-db/planner/lexer"
	"github.com/bengosborn/stream-db/planner/parser"
)

func main() {
	query, err := os.ReadFile("examples/queries/filter_orders.ql")
	if err != nil {
		panic(err)
	}

	tokens, err := lexer.NewTokenizer().Parse(string(query))
	if err != nil {
		panic(err)
	}

	plan, err := parser.NewParser().Parse(tokens)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query:\n%s\n\nQuery plan:\n", query)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(plan); err != nil {
		panic(err)
	}
}
