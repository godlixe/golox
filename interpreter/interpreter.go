package interpreter

import (
	"errors"
	"fmt"
	"golox/ast"
	"golox/statement"
	"golox/token"
)

type GoloxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, argumenst []any) any
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
func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) any {
	return expr.Value
}

// VisitGroupingExpr evaluates a group expression.
func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) any {
	return i.Environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) any {
	value := i.evaluate(expr.Value)
	i.Environment.assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) any {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.Call) any {
	callee := i.evaluate(expr.Callee)

	arguments := []any{}
	for _, argument := range expr.Arguments {
		res := i.evaluate(argument)
		arguments = append(arguments, res)
	}

	if _, ok := callee.(GoloxCallable); !ok {
		// TODO : add runtime error
		fmt.Println("Callee is not a golox callable.", expr.Callee)
		return nil
	}

	var function GoloxCallable = callee.(GoloxCallable)

	if len(arguments) != function.Arity() {
		// TODO : add runtime error
		fmt.Printf("Expected %v arguments but got %v.", function.Arity(), len(arguments))
		return nil
	}

	fnCall := function.Call(i, arguments)
	return fnCall
}

// evaluate evaluates an expression.
func (i *Interpreter) evaluate(expr ast.Expr) any {
	res := expr.Accept(i)
	return res
}

// VisitUnaryExpr evaluates a unary expression.
func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
	right := i.evaluate(expr.Right)

	v, _ := right.(float64)

	switch expr.Operator.Type {
	case token.MINUS:
		err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return err
		}

		return -v
	case token.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

// VisitBinaryExpr evaluates a binary expression.
func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case token.GREATER:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) >= right.(float64)
	case token.LESS:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) <= right.(float64)
	case token.MINUS:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) - right.(float64)
	case token.PLUS:
		if vLeft, ok := left.(string); ok {
			if vRight, ok := right.(string); ok {
				return vLeft + vRight
			}
		} else if vLeft, ok := left.(float64); ok {
			if vRight, ok := right.(float64); ok {
				return vLeft + vRight
			}
		}

		return errors.New("operands must be two numbers or two strings")
	case token.SLASH:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) / right.(float64)
	case token.STAR:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return err
		}
		return left.(float64) * right.(float64)
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *statement.Expression) any {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt *statement.Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(value)
	return value
}

func (i *Interpreter) VisitVarStmt(stmt *statement.Variable) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.Environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *statement.Block) any {
	return i.ExecuteBlock(stmt.Statements, NewEnvironment(i.Environment))
}

func (i *Interpreter) ExecuteBlock(statements []statement.Stmt, environment Environment) any {
	previous := i.Environment

	i.Environment = environment

	var res any = nil

	for _, statement := range statements {
		res = i.execute(statement)
	}

	i.Environment = previous
	return res
}

func (i *Interpreter) VisitIfStmt(stmt *statement.If) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *statement.While) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *statement.Function) any {
	fun := &GoloxFunction{
		Declaration: *stmt,
	}

	i.Environment.Define(stmt.Name.Lexeme, fun)

	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *statement.Return) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	return value
}

func (i *Interpreter) execute(stmt statement.Stmt) any {
	return stmt.Accept(i)
}

// Interpret interprets expressions from an AST.
func (i *Interpreter) Interpret(statements []statement.Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}
