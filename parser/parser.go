package parser

import (
	"errors"
	"fmt"
	"golox/ast"
	errorx "golox/error"
	"golox/statement"
	"golox/token"
)

/*
Below are the current production rules
for golox

expression     → assignment ;

assignment     → IDENTIFIER "=" assignment | logic_or;

logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
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

// Peek returns the current token.
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
	return p.assignment()
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

// finishCall parses each of the arguments to a function
// and includes it in the call node.
func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		if len(arguments) >= 255 {
			return nil, errors.New("Can't have more than 255 arguments.")
		}

		// get first argument
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, expr)

		// get next arguments
		for p.match(token.COMMA) {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

// call parses a function call, determines the callee, and
// calls finishCall() to construct the nodes for a
// function call.
func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
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

	return p.call()
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
			return nil, err
		}

		return &ast.Grouping{
			Expression: expr,
		}, err
	}

	return &ast.Binary{}, nil
}

// ifStatement parses an if statement.
func (p *Parser) ifStatement() (statement.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after 'if'.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch statement.Stmt = nil
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &statement.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

// block parses a block statement. A block statement is
// a set of statements that is enclosed in curly brackets "{}".
func (p *Parser) block() ([]statement.Stmt, error) {
	var statements []statement.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, declaration)
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

// printStatement parses a print statement.
func (p *Parser) printStatement() (statement.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		fmt.Println(value)
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &statement.Print{
		Expression: value,
	}, nil
}

// expressionStatement parses expression statements. Expression
// statements are statements that produces values.
func (p *Parser) expressionStatement() (statement.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &statement.Expression{
		Expression: value,
	}, nil
}

// whileStatement parses a while statement.
func (p *Parser) whileStatement() (statement.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &statement.While{
		Condition: condition,
		Body:      body,
	}, nil
}

// forStatement parses a for statement.
func (p *Parser) forStatement() (statement.Stmt, error) {
	var err error
	var initializer statement.Stmt

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr = nil
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
		if err != nil {
			return nil, err
		}
	}

	var increment ast.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
		if err != nil {
			return nil, err
		}
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &statement.Block{
			Statements: []statement.Stmt{
				body,
				&statement.Expression{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{
			Value: true,
		}
	}

	body = &statement.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &statement.Block{
			Statements: []statement.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

// statement parses statements.
func (p *Parser) statement() (statement.Stmt, error) {
	if p.match(token.FOR) {
		return p.forStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.RETURN) {
		return p.returnStatement()
	}

	if p.match(token.WHILE) {
		return p.whileStatement()
	}

	if p.match(token.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}

		return &statement.Block{
			Statements: statements,
		}, nil
	}

	return p.expressionStatement()
}

// and parses and expressions.
func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// or parses or expressions.
func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, err
}

// assignment parses assignment expressions.
func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*ast.Variable); ok {
			name := v.Name
			return &ast.Assign{
				Name:  name,
				Value: value,
			}, nil
		}

		return nil, errors.New(fmt.Sprint(equals) + "Invalid assignment target. ")
	}

	return expr, nil
}

// varDeclaration parses variable declarations.
func (p *Parser) varDeclaration() (statement.Stmt, error) {
	var (
		initializer ast.Expr
		err         error
	)

	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return &statement.Variable{
		Name:        name,
		Initializer: initializer,
	}, nil
}

// function parses functions.
func (p *Parser) function(kind string) (*statement.Function, error) {
	name, err := p.consume(token.IDENTIFIER, fmt.Sprintf("Expect %v name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %v name.", kind))
	if err != nil {
		return nil, err
	}

	var parameters []token.Token
	if !p.check(token.RIGHT_PAREN) {
		if len(parameters) >= 255 {
			return nil, errors.New("Can't have more than 255 parameters.")
		}

		param, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, param)

		for p.match(token.COMMA) {
			param, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, param)
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %v body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &statement.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

// declaration parses declarations.
func (p *Parser) declaration() (statement.Stmt, error) {
	if p.match(token.FUN) {
		return p.function("function")
	}

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

// returnStatement parses return statements.
func (p *Parser) returnStatement() (statement.Stmt, error) {
	keyword := p.previous()
	var value ast.Expr = nil
	var err error
	if !p.check(token.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &statement.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}

// parse parses the tokens inside the token list.
func (p *Parser) Parse() ([]statement.Stmt, bool) {

	var statements []statement.Stmt
	var isError bool

	for !p.isAtEnd() {
		statement, err := p.declaration()

		if err != nil {
			isError = true
			errorx.Error(p.peek().Line, err.Error())
			p.synchronize()
		}

		statements = append(statements, statement)
	}

	return statements, isError
}
