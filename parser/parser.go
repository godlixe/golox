package parser

import (
	"errors"
	"fmt"
	"golox/ast"
	"golox/statement"
	"golox/token"
)

/*
Below are the current production rules
for golox

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/

// Parser represents a parser object.
// Tokens contains the list of tokens scanned.
// Current is the current token being parsed from the list.
type Parser struct {
	Tokens  []token.Token
	Current int
}

// Previous returns the previous token.
func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}

// Peek returnsthe current token.
func (p *Parser) peek() token.Token {
	return p.Tokens[p.Current]
}

// isAtEnd checks if the current token is
// an EOF-token.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// advance increments the current pointer from
// the parser object if the parser is not at the
// end of the token list and return the previous token
// (which should be the token before the increment).
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.Current += 1
	}

	return p.previous()
}

// check checks if the current token's type is equal to tp.
func (p *Parser) check(tp token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tp
}

// match checks the current token's type with given types and
// returns true if it matches one.
func (p *Parser) match(types ...token.TokenType) bool {
	for _, tp := range types {
		if p.check(tp) {
			p.advance()
			return true
		}
	}

	return false
}

// consume checks if the current token is of type tp and returns
// error if the token does not match the type.
func (p *Parser) consume(tp token.TokenType, message string) (token.Token, error) {
	if p.check(tp) {
		return p.advance(), nil
	}

	return p.peek(), errors.New(message)
}

// synchronize unwinds the parser by discarding tokens.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		p.advance()
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}

		p.advance()
	}
}

// expression parses an expression starting from the
// expression of lowest precedence.
func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

// equality parses an equality expression. An equality
// expression contains a comparison with != or ==.
func (p *Parser) equality() (ast.Expr, error) {
	var err error
	expr, err := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// comparison parses a comparison expression. A comparison
// expression contains a comparison with >, >=, <, <=.
func (p *Parser) comparison() (ast.Expr, error) {
	var err error
	expr, err := p.term()

	for p.match(
		token.GREATER,
		token.GREATER_EQUAL,
		token.LESS,
		token.LESS_EQUAL,
	) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// comparison parses a term expression. A term
// expression contains addition or subtraction.
func (p *Parser) term() (ast.Expr, error) {
	var err error
	expr, err := p.factor()

	for p.match(
		token.PLUS,
		token.MINUS,
	) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// factor parses a factor expression. A factor
// expression contains multiplication and division.
func (p *Parser) factor() (ast.Expr, error) {
	var err error
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(
		token.STAR,
		token.SLASH,
	) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// unary parses a unary expression. A unary
// expression contains negation (! or -).
func (p *Parser) unary() (ast.Expr, error) {
	for p.match(
		token.BANG,
		token.MINUS,
	) {
		operator := p.previous()
		right, err := p.unary()
		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, err
	}

	return p.primary()
}

// primary parses a primary expression. A primary
// expression contains booleans, nil, numbers, strings, and
// expressions inside parentheses.
func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{
			Value: false,
		}, nil
	}

	if p.match(token.TRUE) {
		return &ast.Literal{
			Value: true,
		}, nil
	}

	if p.match(token.NIL) {
		return &ast.Literal{
			Value: nil,
		}, nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(token.IDENTIFIER) {
		return &ast.Variable{
			Name: p.previous(),
		}, nil
	}

	if p.match(token.LEFT_PAREN) {
		var err error
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			fmt.Println(err)
		}

		return &ast.Grouping{
			Expression: expr,
		}, err
	}

	return &ast.Binary{}, nil
}

func (p *Parser) printStatement() (statement.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		fmt.Println(value)
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &statement.Print{
		Expression: value,
	}, nil
}

func (p *Parser) expressionStatement() (statement.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		fmt.Println(value)
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &statement.Expression{
		Expression: value,
	}, nil
}

func (p *Parser) statement() (statement.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) varDeclaration() (statement.Stmt, error) {
	var (
		initializer ast.Expr
		err         error
	)

	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration")
	return &statement.Variable{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) declaration() (statement.Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

// parse parses the tokens inside the token list.
func (p *Parser) Parse() []statement.Stmt {

	var statements []statement.Stmt

	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			fmt.Println(err)
			p.synchronize()
		}

		statements = append(statements, statement)
	}

	return statements
}
