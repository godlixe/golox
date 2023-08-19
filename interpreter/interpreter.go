package interpreter

import (
	"errors"
	"fmt"
	"golox/ast"
	"golox/token"
)

type Interpreter struct {
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

// TODO
func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

// VisitUnaryExpr evaluates a unary expression.
func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
	right := i.evaluate(expr.Right)

	v, ok := right.(float64)
	if !ok {
		fmt.Println("Type assertion error")
	}

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

// Interpret interprets expressions from an AST.
func (i *Interpreter) Interpret(expression ast.Expr) {
	value := i.evaluate(expression)
	fmt.Println(value)
}
