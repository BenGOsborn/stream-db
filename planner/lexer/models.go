package lexer

type TokenType uint8

const (
	TokenTypeSelect         TokenType = 1
	TokenTypeIdentifier     TokenType = 2
	TokenTypeComma          TokenType = 3
	TokenTypeParenthesis    TokenType = 4
	TokenTypeBinaryOperator TokenType = 5
	TokenTypeSemicolon      TokenType = 6
	TokenTypeFrom           TokenType = 7
	TokenTypeWhere          TokenType = 8
	TokenTypeCondition      TokenType = 9
)

type Token struct {
	Type  TokenType
	Value string
}
