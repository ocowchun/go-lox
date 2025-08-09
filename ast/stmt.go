package ast

import (
	"github.com/ocowchun/go-lox/token"
)

type Stmt interface {
	Stmt()
	Accept(visitor StmtVisitor) any
}

type StmtVisitor interface {
	VisitExpressionStatement(stmt *ExpressionStatement) any
	VisitPrintStatement(stmt *PrintStatement) any
	VisitVarStatement(stmt *VarStatement) any
	VisitBlockStatement(stmt *BlockStatement) any
	VisitIfStatement(stmt *IfStatement) any
	VisitWhileStatement(stmt *WhileStatement) any
	VisitFunctionStatement(stmt *FunctionStatement) any
	VisitReturnStatement(stmt *ReturnStatement) any
	VisitClassStatement(stmt *ClassStatement) any
}

type ExpressionStatement struct {
	Expression Expr
}

func (stmt *ExpressionStatement) Stmt() {}

func (stmt *ExpressionStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpressionStatement(stmt)
}

type PrintStatement struct {
	Expression Expr
}

func (stmt *PrintStatement) Stmt() {}

func (stmt *PrintStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStatement(stmt)
}

type VarStatement struct {
	Name        token.Token
	Initializer Expr
}

func (stmt *VarStatement) Stmt() {}

func (stmt *VarStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStatement(stmt)
}

type BlockStatement struct {
	Statements []Stmt
}

func (stmt *BlockStatement) Stmt() {}

func (stmt *BlockStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStatement(stmt)
}

type IfStatement struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (stmt *IfStatement) Stmt() {}

func (stmt *IfStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitIfStatement(stmt)
}

type WhileStatement struct {
	Condition Expr
	Body      Stmt
}

func (stm *WhileStatement) Stmt() {}

func (stm *WhileStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhileStatement(stm)
}

type FunctionStatement struct {
	Name       token.Token
	Parameters []token.Token
	Body       *BlockStatement
}

func (stmt *FunctionStatement) Stmt() {}

func (stmt *FunctionStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunctionStatement(stmt)
}

type ReturnStatement struct {
	// keep Keyword, so we can use its location for error reporting
	Keyword token.Token
	Value   Expr
}

func (stmt *ReturnStatement) Stmt() {}

func (stmt *ReturnStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitReturnStatement(stmt)
}

type ClassStatement struct {
	Name token.Token
	// nil if no superclass
	Superclass *VariableExpression
	Methods    []*FunctionStatement
}

func (stmt *ClassStatement) Stmt() {}

func (stmt *ClassStatement) Accept(visitor StmtVisitor) any {
	return visitor.VisitClassStatement(stmt)
}
