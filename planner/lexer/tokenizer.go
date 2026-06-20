package lexer

type Tokenizer interface {
	Parse(raw string) ([]Token, error)
}

type tokenizer struct{}

type tokenParseResult struct {
	Token     Token
	NextIndex int
	Matched   bool
}

type tokenParser func(string, int) tokenParseResult

var tokenParsers = []tokenParser{
	parseSelect,
	parseFrom,
	parseWhere,
	parseCondition,
	parseComma,
	parseParenthesis,
	parseBinaryOperator,
	parseSemicolon,
	parseIdentifier,
}

func NewTokenizer() Tokenizer {
	return &tokenizer{}
}

func (t *tokenizer) Parse(raw string) ([]Token, error) {
	tokens := make([]Token, 0)

	for i := 0; i < len(raw); {
		matched := false

		for _, parser := range tokenParsers {
			result := parser(raw, i)
			if !result.Matched {
				continue
			}

			tokens = append(tokens, result.Token)
			i = result.NextIndex
			matched = true
			break
		}

		if !matched {
			i += 1
		}
	}

	return tokens, nil
}

func parseSelect(raw string, index int) tokenParseResult {
	return parseKeyword(raw, index, "SELECT", TokenTypeSelect)
}

func parseFrom(raw string, index int) tokenParseResult {
	return parseKeyword(raw, index, "FROM", TokenTypeFrom)
}

func parseWhere(raw string, index int) tokenParseResult {
	return parseKeyword(raw, index, "WHERE", TokenTypeWhere)
}

func parseCondition(raw string, index int) tokenParseResult {
	for _, condition := range []string{"AND", "OR"} {
		result := parseKeyword(raw, index, condition, TokenTypeCondition)
		if result.Matched {
			return result
		}
	}

	return tokenParseResult{}
}

func parseComma(raw string, index int) tokenParseResult {
	return parseCharacter(raw, index, ',', TokenTypeComma)
}

func parseParenthesis(raw string, index int) tokenParseResult {
	if index >= len(raw) || raw[index] != '(' && raw[index] != ')' {
		return tokenParseResult{}
	}

	return tokenParseResult{
		Token: Token{
			Type:  TokenTypeParenthesis,
			Value: raw[index : index+1],
		},
		NextIndex: index + 1,
		Matched:   true,
	}
}

func parseBinaryOperator(raw string, index int) tokenParseResult {
	if index >= len(raw) || raw[index] != '=' && raw[index] != '<' && raw[index] != '>' {
		return tokenParseResult{}
	}

	return tokenParseResult{
		Token: Token{
			Type:  TokenTypeBinaryOperator,
			Value: raw[index : index+1],
		},
		NextIndex: index + 1,
		Matched:   true,
	}
}

func parseSemicolon(raw string, index int) tokenParseResult {
	return parseCharacter(raw, index, ';', TokenTypeSemicolon)
}

func parseIdentifier(raw string, index int) tokenParseResult {
	if index >= len(raw) {
		return tokenParseResult{}
	}

	if raw[index] == '\'' {
		nextIndex := index + 1
		for nextIndex < len(raw) && raw[nextIndex] != '\'' {
			nextIndex++
		}

		if nextIndex < len(raw) {
			nextIndex++
		}

		return tokenParseResult{
			Token: Token{
				Type:  TokenTypeIdentifier,
				Value: raw[index:nextIndex],
			},
			NextIndex: nextIndex,
			Matched:   true,
		}
	}

	if !isIdentifierCharacter(raw[index]) {
		return tokenParseResult{}
	}

	nextIndex := index + 1
	for nextIndex < len(raw) && isIdentifierCharacter(raw[nextIndex]) {
		nextIndex++
	}

	return tokenParseResult{
		Token: Token{
			Type:  TokenTypeIdentifier,
			Value: raw[index:nextIndex],
		},
		NextIndex: nextIndex,
		Matched:   true,
	}
}

func parseKeyword(raw string, index int, keyword string, tokenType TokenType) tokenParseResult {
	nextIndex := index + len(keyword)
	if nextIndex > len(raw) || raw[index:nextIndex] != keyword {
		return tokenParseResult{}
	}

	if nextIndex < len(raw) && isIdentifierCharacter(raw[nextIndex]) {
		return tokenParseResult{}
	}

	return tokenParseResult{
		Token: Token{
			Type:  tokenType,
			Value: raw[index:nextIndex],
		},
		NextIndex: nextIndex,
		Matched:   true,
	}
}

func parseCharacter(raw string, index int, character byte, tokenType TokenType) tokenParseResult {
	if index >= len(raw) || raw[index] != character {
		return tokenParseResult{}
	}

	return tokenParseResult{
		Token: Token{
			Type:  tokenType,
			Value: raw[index : index+1],
		},
		NextIndex: index + 1,
		Matched:   true,
	}
}

func isIdentifierCharacter(character byte) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' ||
		character == '_'
}
