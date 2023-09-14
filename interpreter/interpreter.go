package interpreter

import (
	"errors"
	"fmt"
	"golox/ast"
	"golox/statement"
	"golox/token"
	"os"
)

type GoloxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, argumenst []any) (any, error)
}

type Interpreter struct {
	Environment Environment
	Globals     Environment
}

// isTruthy checks if an object is truthy or falsey.
// everything except nil and false is considered truthy,
// else is falsey.
func (u *Interpreter) isTruthy(obj any) bool {
	switch v := obj.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

// isEqual checks if two objects are equal.
func (i *Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

// checkNumberOperand checks if an operand is a number.
// a number is defined as numbers that can be parsed to Go's float64.
func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}

	return errors.New("operand must be a number")
}

// checkNumberOperands checks if two operands are numbers.
// a number is defined as numbers that can be parsed to Go's float64.
func (i *Interpreter) checkNumberOperands(operator token.Token, left any, right any) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}

	return errors.New("operands must be numbers")
}

// VisitLiteralExpr evaluates literal expression.
func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

// VisitGroupingExpr evaluates a group expression.
func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) (any, error) {
	return i.Environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	i.Environment.assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.Call) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argument := range expr.Arguments {
		res, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, res)
	}

	if _, ok := callee.(GoloxCallable); !ok {
		return nil, fmt.Errorf("Callee is not a golox callable.", expr.Callee)
	}

	var function GoloxCallable = callee.(GoloxCallable)

	if len(arguments) != function.Arity() {
		return nil, fmt.Errorf("Expected %v arguments but got %v.", function.Arity(), len(arguments))
	}

	fnCall, err := function.Call(i, arguments)
	if err != nil {
		return nil, err
	}

	return fnCall, nil
}

// evaluate evaluates an expression.
func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {

	// error is to detect runtime errors
	res, err := expr.Accept(i)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// VisitUnaryExpr evaluates a unary expression.
func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	v, _ := right.(float64)

	switch expr.Operator.Type {
	case token.MINUS:
		err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return -v, nil
	case token.BANG:
		return !i.isTruthy(right), nil
	}

	return nil, nil
}

// VisitBinaryExpr evaluates a binary expression.
func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case token.GREATER:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.MINUS:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.PLUS:
		if vLeft, ok := left.(string); ok {
			if vRight, ok := right.(string); ok {
				return vLeft + vRight, nil
			}
		} else if vLeft, ok := left.(float64); ok {
			if vRight, ok := right.(float64); ok {
				return vLeft + vRight, nil
			}
		}

		return nil, errors.New("operands must be two numbers or two strings")
	case token.SLASH:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.STAR:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *statement.Expression) (any, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt *statement.Print) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(value)
	return value, nil
}

func (i *Interpreter) VisitVarStmt(stmt *statement.Variable) (any, error) {
	var value any
	var err error

	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.Environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *statement.Block) (any, error) {
	return i.ExecuteBlock(stmt.Statements, NewEnvironment(i.Environment))
}

func (i *Interpreter) ExecuteBlock(statements []statement.Stmt, environment Environment) (any, error) {
	previous := i.Environment

	i.Environment = environment

	var res any = nil
	var err error

	for _, statement := range statements {
		res, err = i.execute(statement)
		if err != nil {
			return nil, err
		}
	}

	i.Environment = previous
	return res, nil
}

func (i *Interpreter) VisitIfStmt(stmt *statement.If) (any, error) {
	res, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(res) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *statement.While) (any, error) {
	res, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	for i.isTruthy(res) {
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}

		res, err = i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *statement.Function) (any, error) {
	fun := &GoloxFunction{
		Declaration: *stmt,
	}

	i.Environment.Define(stmt.Name.Lexeme, fun)

	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *statement.Return) (any, error) {
	var value any = nil
	var err error

	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (i *Interpreter) execute(stmt statement.Stmt) (any, error) {
	return stmt.Accept(i)
}

// Interpret interprets expressions from an AST.
func (i *Interpreter) Interpret(statements []statement.Stmt) {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
