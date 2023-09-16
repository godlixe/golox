package scanner

import (
	"golox/token"
	"testing"
)

func TestToken(t *testing.T) {
	input := `( ) { } , . - + ; * ! != == = <= < >= > / test //end 
	"" $`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.MINUS, "-"},
		{token.PLUS, "+"},
		{token.SEMICOLON, ";"},
		{token.STAR, "*"},
		{token.BANG, "!"},
		{token.BANG_EQUAL, "!="},
		{token.EQUAL_EQUAL, "=="},
		{token.EQUAL, "="},
		{token.LESS_EQUAL, "<="},
		{token.LESS, "<"},
		{token.GREATER_EQUAL, ">="},
		{token.GREATER, ">"},
		{token.SLASH, "/"},
		{token.IDENTIFIER, "test"},
	}

	scanner := New(input)

	tokens := scanner.ScanTokens()
	for i, tt := range tests {

		if tokens[i].Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%v, got=%v", i, tt.expectedType, tokens[i].Type)
		}

		if tokens[i].Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%v, got=%v", i, tt.expectedLiteral, tokens[i].Literal)
		}
	}
}
