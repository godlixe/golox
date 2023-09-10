package statement

import (
	"golox/ast"
	"golox/token"
)

/*
Statement production rules

program        → statement* EOF ;

declaration    → funDecl
			   | varDecl
			   | statement ;

funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;

parameters     → IDENTIFIER ( "," IDENTIFIER )* ;

statement      → exprStmt
			   | forStmt
			   | ifStmt
               | printStmt
			   | returnStmt
			   | whileStmt
			   | block ;

returnStmt      → "return" expression? ";" ;

forStmt        → "for" "(" (varDecl | exprStmt | ";")
				expression? ";"
				expression? ")" statement ;
whileStmt      → "while" "(" expression ")" statement ;
ifStmt		   → "if" "(" expression ")" statement
				("else" statement)? ;
block          → "{" declaration* "}";
exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
*/

type Visitor interface {
	VisitBlockStmt(stmt *Block) any
	// VisitClassStmt(stmt *Class)
	VisitExpressionStmt(stmt *Expression) any
	VisitFunctionStmt(stmt *Function) any
	VisitIfStmt(stmt *If) any
	VisitPrintStmt(stmt *Print) any
	VisitReturnStmt(stmt *Return) any
	VisitVarStmt(stmt *Variable) any
	VisitWhileStmt(stmt *While) any
}

type Stmt interface {
	Accept(visitor Visitor) any
}

type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(visitor Visitor) any {
	return visitor.VisitBlockStmt(b)
}

type Class struct {
	Name       token.Token
	SuperClass ast.Variable
	Methods    []Function
}

// func (c *Class) Accept(visitor Visitor) {
// 	visitor.VisitClassStmt(c)
// }

type Expression struct {
	Expression ast.Expr
}

func (e *Expression) Accept(visitor Visitor) any {
	return visitor.VisitExpressionStmt(e)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (f *Function) Accept(visitor Visitor) any {
	return visitor.VisitFunctionStmt(f)
}

type If struct {
	Condition  ast.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(visitor Visitor) any {
	return visitor.VisitIfStmt(i)
}

type Print struct {
	Expression ast.Expr
}

func (p *Print) Accept(visitor Visitor) any {
	return visitor.VisitPrintStmt(p)
}

type Return struct {
	Keyword token.Token
	Value   ast.Expr
}

func (r *Return) Accept(visitor Visitor) any {
	return visitor.VisitReturnStmt(r)
}

type Variable struct {
	Name        token.Token
	Initializer ast.Expr
}

func (v *Variable) Accept(visitor Visitor) any {
	return visitor.VisitVarStmt(v)
}

type While struct {
	Condition ast.Expr
	Body      Stmt
}

func (w *While) Accept(visitor Visitor) any {
	return visitor.VisitWhileStmt(w)
}
