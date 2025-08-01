package ast

import (
	"fmt"
	"github.com/ocowchun/go-lox/token"
	"strings"
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
	Name        *token.Token
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

type StatementPrinter struct {
	expressionPrinter *ExpressionPrinter
}

func NewStatementPrinter() *StatementPrinter {
	return &StatementPrinter{
		expressionPrinter: &ExpressionPrinter{},
	}

}

func (printer *StatementPrinter) Print(stmt Stmt) string {
	res := stmt.Accept(printer).(string)
	return res
}

func (printer *StatementPrinter) VisitExpressionStatement(stmt *ExpressionStatement) any {
	return stmt.Expression.Accept(printer.expressionPrinter)
}

func (printer *StatementPrinter) VisitPrintStatement(stmt *PrintStatement) any {
	return fmt.Sprintf("(print %s)", stmt.Expression.Accept(printer.expressionPrinter))
}

func (printer *StatementPrinter) VisitVarStatement(stmt *VarStatement) any {
	return fmt.Sprintf("(define %s %s)", stmt.Name.Lexeme, stmt.Initializer.Accept(printer.expressionPrinter))
}

func (printer *StatementPrinter) VisitBlockStatement(stmt *BlockStatement) any {
	var b strings.Builder
	b.WriteString("(begin\n")
	for _, s := range stmt.Statements {
		b.WriteString(printer.Print(s))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}
