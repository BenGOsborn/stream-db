package parser

import (
	"fmt"

	"github.com/bengosborn/stream-db/planner/lexer"
)

type Parser interface {
	Parse(tokens []lexer.Token) (QueryPlan, error)
}

type parser struct{}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(tokens []lexer.Token) (QueryPlan, error) {
	index := 0

	selectAction, index, err := parseSelect(tokens, index)
	if err != nil {
		return QueryPlan{}, err
	}

	fromAction, index, err := parseFrom(tokens, index)
	if err != nil {
		return QueryPlan{}, err
	}

	where, index, err := parseExpression(tokens, index)
	if err != nil {
		return QueryPlan{}, err
	}

	_, index, err = consume(tokens, index, lexer.TokenTypeSemicolon)
	if err != nil {
		return QueryPlan{}, err
	}

	if index != len(tokens) {
		return QueryPlan{}, fmt.Errorf("unexpected token %q at index %d", tokens[index].Value, index)
	}

	return QueryPlan{
		Select: selectAction,
		From:   fromAction,
		Where:  where,
	}, nil
}

func parseSelect(tokens []lexer.Token, index int) (SelectAction, int, error) {
	_, nextIndex, err := consume(tokens, index, lexer.TokenTypeSelect)
	if err != nil {
		return SelectAction{}, index, fmt.Errorf("parse select: %w", err)
	}

	columns, nextIndex, err := parseIdentifierList(tokens, nextIndex, lexer.TokenTypeFrom)
	if err != nil {
		return SelectAction{}, index, fmt.Errorf("parse select: %w", err)
	}

	return SelectAction{Columns: columns}, nextIndex, nil
}

func parseFrom(tokens []lexer.Token, index int) (FromAction, int, error) {
	_, nextIndex, err := consume(tokens, index, lexer.TokenTypeFrom)
	if err != nil {
		return FromAction{}, index, fmt.Errorf("parse from: %w", err)
	}

	tables, nextIndex, err := parseIdentifierList(
		tokens,
		nextIndex,
		lexer.TokenTypeWhere,
		lexer.TokenTypeSemicolon,
	)
	if err != nil {
		return FromAction{}, index, fmt.Errorf("parse from: %w", err)
	}

	return FromAction{Tables: tables}, nextIndex, nil
}

func parseIdentifierList(
	tokens []lexer.Token,
	index int,
	endTypes ...lexer.TokenType,
) ([]string, int, error) {
	values := make([]string, 0)
	nextIndex := index
	expectIdentifier := true

	for !atEnd(tokens, nextIndex) && !matchesAny(tokens, nextIndex, endTypes...) {
		if expectIdentifier {
			token, parsedIndex, err := consume(tokens, nextIndex, lexer.TokenTypeIdentifier)
			if err != nil {
				return nil, index, err
			}

			values = append(values, token.Value)
			nextIndex = parsedIndex
			expectIdentifier = false
			continue
		}

		_, parsedIndex, err := consume(tokens, nextIndex, lexer.TokenTypeComma)
		if err != nil {
			return nil, index, err
		}

		nextIndex = parsedIndex
		expectIdentifier = true
	}

	if len(values) == 0 {
		return nil, index, fmt.Errorf("expected at least one identifier")
	}

	if expectIdentifier {
		return nil, index, fmt.Errorf("expected identifier after comma")
	}

	return values, nextIndex, nil
}

func parseExpression(tokens []lexer.Token, index int) (*Filter, int, error) {
	if !matches(tokens, index, lexer.TokenTypeWhere) {
		return nil, index, nil
	}

	filter, nextIndex, err := parseOr(tokens, index+1)
	if err != nil {
		return nil, index, fmt.Errorf("parse where: %w", err)
	}

	return filter, nextIndex, nil
}

func parseOr(tokens []lexer.Token, index int) (*Filter, int, error) {
	left, nextIndex, err := parseAnd(tokens, index)
	if err != nil {
		return nil, index, err
	}

	for matchesValue(tokens, nextIndex, lexer.TokenTypeCondition, "OR") {
		right, parsedIndex, err := parseAnd(tokens, nextIndex+1)
		if err != nil {
			return nil, index, err
		}

		left = &Filter{
			Left:            left,
			LogicalOperator: LogicalOperatorOr,
			Right:           right,
		}
		nextIndex = parsedIndex
	}

	return left, nextIndex, nil
}

func parseAnd(tokens []lexer.Token, index int) (*Filter, int, error) {
	left, nextIndex, err := parsePrimary(tokens, index)
	if err != nil {
		return nil, index, err
	}

	for matchesValue(tokens, nextIndex, lexer.TokenTypeCondition, "AND") {
		right, parsedIndex, err := parsePrimary(tokens, nextIndex+1)
		if err != nil {
			return nil, index, err
		}

		left = &Filter{
			Left:            left,
			LogicalOperator: LogicalOperatorAnd,
			Right:           right,
		}
		nextIndex = parsedIndex
	}

	return left, nextIndex, nil
}

func parsePrimary(tokens []lexer.Token, index int) (*Filter, int, error) {
	if matchesValue(tokens, index, lexer.TokenTypeParenthesis, "(") {
		filter, nextIndex, err := parseExpression(tokens, index+1)
		if err != nil {
			return nil, index, err
		}

		_, nextIndex, err = consumeValue(
			tokens,
			nextIndex,
			lexer.TokenTypeParenthesis,
			")",
		)
		if err != nil {
			return nil, index, err
		}

		return filter, nextIndex, nil
	}

	return parseComparison(tokens, index)
}

func parseComparison(tokens []lexer.Token, index int) (*Filter, int, error) {
	left, nextIndex, err := consume(tokens, index, lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, index, err
	}

	operator, nextIndex, err := consume(tokens, nextIndex, lexer.TokenTypeBinaryOperator)
	if err != nil {
		return nil, index, err
	}

	right, nextIndex, err := consume(tokens, nextIndex, lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, index, err
	}

	return &Filter{
		Comparison: &Comparison{
			Left:     left.Value,
			Operator: operator.Value,
			Right:    right.Value,
		},
	}, nextIndex, nil
}

func consume(
	tokens []lexer.Token,
	index int,
	tokenType lexer.TokenType,
) (lexer.Token, int, error) {
	if atEnd(tokens, index) {
		return lexer.Token{}, index, fmt.Errorf("expected token type %d at end of query", tokenType)
	}

	token := tokens[index]
	if token.Type != tokenType {
		return lexer.Token{}, index, fmt.Errorf(
			"expected token type %d at index %d, got %q",
			tokenType,
			index,
			token.Value,
		)
	}

	return token, index + 1, nil
}

func consumeValue(
	tokens []lexer.Token,
	index int,
	tokenType lexer.TokenType,
	value string,
) (lexer.Token, int, error) {
	if !matchesValue(tokens, index, tokenType, value) {
		if atEnd(tokens, index) {
			return lexer.Token{}, index, fmt.Errorf("expected %q at end of query", value)
		}

		return lexer.Token{}, index, fmt.Errorf(
			"expected %q at index %d, got %q",
			value,
			index,
			tokens[index].Value,
		)
	}

	return tokens[index], index + 1, nil
}

func matches(tokens []lexer.Token, index int, tokenType lexer.TokenType) bool {
	return !atEnd(tokens, index) && tokens[index].Type == tokenType
}

func matchesAny(
	tokens []lexer.Token,
	index int,
	tokenTypes ...lexer.TokenType,
) bool {
	for _, tokenType := range tokenTypes {
		if matches(tokens, index, tokenType) {
			return true
		}
	}

	return false
}

func matchesValue(
	tokens []lexer.Token,
	index int,
	tokenType lexer.TokenType,
	value string,
) bool {
	return matches(tokens, index, tokenType) && tokens[index].Value == value
}

func atEnd(tokens []lexer.Token, index int) bool {
	return index >= len(tokens)
}
