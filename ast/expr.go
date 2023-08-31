package ast

import (
	"fmt"
	"golox/token"
)

// Visitor is an interface for implementing
// the visitor pattern. It provides method
// for visiting expressions in the AST.
type Visitor interface {
	VisitAssignExpr(assign *Assign) any
	VisitBinaryExpr(binary *Binary) any
	// VisitCallExpr(call *Call) any
	// VisitGetExpr(get *Get) any
	VisitGroupingExpr(grouping *Grouping) any
	VisitLiteralExpr(literal *Literal) any
	VisitLogicalExpr(logical *Logical) any
	// VisitSetExpr(set *Set) any
	// VisitSuperExpr(super *Super) any
	// VisitThisExpr(this *This) any
	VisitUnaryExpr(unary *Unary) any
	VisitVariableExpr(variable *Variable) any
}

func PrintAst(expr Expr) {
	if expr == nil {
		return
	}

	if v, ok := expr.(*Literal); ok {
		fmt.Println(v)
	}

	if v, ok := expr.(*Binary); ok {
		PrintAst(v.Left)
		PrintAst(v.Right)
		fmt.Println(v.Operator)
	}
}

// Expr defines an interface for an expression.
// an Expression could be of assignment, binary
// call, get, grouping, literal, logical, set
// super, this, unary, and variable.
type Expr interface {
	Accept(visitor Visitor) any
}

// Assign represents an assignment
// or a variable declaration.
type Assign struct {
	Name  token.Token
	Value Expr
	Expr
}

func (a *Assign) Accept(visitor Visitor) any {
	return visitor.VisitAssignExpr(a)
}

// Binary represents a binary
// operation.
type Binary struct {
	Left     Expr
	Right    Expr
	Operator token.Token
}

func (b *Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(b)
}

// Call represents a function call.
type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

// func (c *Call) Accept(visitor Visitor) any {
// 	return visitor.VisitCallExpr(c)
// }

// Get represents getting an object's property.
type Get struct {
	Object Expr
	Name   token.Token
}

// func (g *Get) Accept(visitor Visitor) any {
// 	return visitor.VisitGetExpr(g)
// }

// Group represents grouping of expression
// with parentheses.
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(g)
}

// Literal represents literals.
type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(l)
}

// Logical represents logical expressions.
type Logical struct {
	Left     Expr
	Right    Expr
	Operator token.Token
}

func (l *Logical) Accept(visitor Visitor) any {
	return visitor.VisitLogicalExpr(l)
}

// Set sets an object's property to a value.
type Set struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

// func (s *Set) Accept(visitor Visitor) any {
// 	return visitor.VisitSetExpr(s)
// }

// Super represents a superclass.
type Super struct {
	Keyword token.Token
	Method  token.Token
}

// func (s *Super) Accept(visitor Visitor) any {
// 	return visitor.VisitSuperExpr(s)
// }

// This represents a class's self reference.
type This struct {
	Keyword token.Token
}

// func (t *This) Accept(visitor Visitor) any {
// 	return visitor.VisitThisExpr(t)
// }

// Unary represents a unary expression.
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(u)
}

// Variable represents a variable.
type Variable struct {
	Name token.Token
}

func (v *Variable) Accept(visitor Visitor) any {
	return visitor.VisitVariableExpr(v)
}
