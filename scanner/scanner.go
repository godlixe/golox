package scanner

import (
	errorx "golox/error"
	"golox/token"
	"strconv"
)

// keywords contain reserved keywords for the
// golox language.
var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

// Scanner defines a scanner object.
type Scanner struct {
	Start   int
	Current int
	Line    int

	Source string
	Tokens []token.Token
}

// New creates a new Scanner instance.
func New(source string) Scanner {
	return Scanner{
		Source:  source,
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

// isAtEnd checks if the current pointer
// points at the end of the source.
func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

// advance returns the current string
// and increments the current pointer.
func (s *Scanner) advance() string {
	s.Current += 1
	return string(s.Source[s.Current-1])
}

// addToken adds a token to the token list.
func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	text := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.Line,
	})
}

// match matches the current string with
// a string passed to the parameter.
func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}

	if string(s.Source[s.Current]) != expected {
		return false
	}

	s.Current += 1
	return true
}

// peek reads the current character without
// consuming it.
func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return "\\0"
	}

	return string(s.Source[s.Current])
}

// peek reads the next character without
// consuming it.
func (s *Scanner) peekNext() string {
	if s.Current+1 >= len(s.Source) {
		return "\\0"
	}

	return string(s.Source[s.Current+1])
}

// isDigit checks if a single character
// string is a number (0-9).
func isDigit(s string) bool {
	if len(s) > 1 {
		return false
	}

	return s >= "0" && s <= "9"
}

// isAlpha checks if a single character
// string is an alphabet (a-zA-Z_).
func isAlpha(s string) bool {
	return (s >= "a" && s <= "z") ||
		(s >= "A" && s <= "Z") ||
		s == "_"
}

// isAlphaNumeric checks if a single character
// string is an alphabet (a-zA-Z_) or a number (0-9).
func isAlphaNumeric(s string) bool {
	return isAlpha(s) || isDigit(s)
}

// string scans for a string and
// adds it to the token list.
func (s *Scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errorx.Error(s.Line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.Source[s.Start+1 : s.Current-1]
	s.addToken(token.STRING, string(value))
}

// number scans for a number and
// adds it to the token list.
func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && isDigit(s.peekNext()) {

		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(string(s.Source[s.Start:s.Current]), 64)
	if err != nil {
		errorx.Error(s.Line, "Unparsable float")
	}

	s.addToken(token.NUMBER, num)
}

// identifier identifies a reserved keyword
// and adds it to the token list.
func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]
	tokenType := keywords[text]

	if tokenType == "" {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType, text)
}

// scanToken scans a token, matches it
// and adds it to the token list.
func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case "(":
		s.addToken(token.LEFT_PAREN, "(")
	case ")":
		s.addToken(token.RIGHT_PAREN, ")")
	case "{":
		s.addToken(token.LEFT_BRACE, "{")
	case "}":
		s.addToken(token.RIGHT_BRACE, "}")
	case ",":
		s.addToken(token.COMMA, ",")
	case ".":
		s.addToken(token.DOT, ".")
	case "-":
		s.addToken(token.MINUS, "-")
	case "+":
		s.addToken(token.PLUS, "+")
	case ";":
		s.addToken(token.SEMICOLON, ";")
	case "*":
		s.addToken(token.STAR, "*")
	case "!":
		if s.match("=") {
			s.addToken(token.BANG_EQUAL, "!=")
		} else {
			s.addToken(token.BANG, "!")
		}
	case "=":
		if s.match("=") {
			s.addToken(token.EQUAL_EQUAL, "==")
		} else {
			s.addToken(token.EQUAL, "=")
		}
	case "<":
		if s.match("=") {
			s.addToken(token.LESS_EQUAL, "<=")
		} else {
			s.addToken(token.LESS, "<")
		}
	case ">":
		if s.match("=") {
			s.addToken(token.GREATER_EQUAL, ">=")
		} else {
			s.addToken(token.GREATER, ">")
		}
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, "/")
		}
	case " ":
	case "\r":
	case "\t":
	case "\n":
		s.Line++
	case "\"":
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			errorx.Error(s.Line, "Unexpected character "+c)
		}
	}
}

// scanTokens scans the source from start to end,
// adding an EOF token to the end of the token list.
func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.Start = s.Current
		s.scanToken()
	}

	s.Tokens = append(
		s.Tokens,
		token.Token{
			Type:    token.EOF,
			Lexeme:  "",
			Literal: nil,
			Line:    s.Line,
		},
	)

	return s.Tokens
}
