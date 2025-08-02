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
	VisitIfStatement(stmt *IfStatement) any
	VisitWhileStatement(stmt *WhileStatement) any
	VisitFunctionStatement(stmt *FunctionStatement) any
	VisitReturnStatement(stmt *ReturnStatement) any
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

func (printer *StatementPrinter) VisitIfStatement(stmt *IfStatement) any {
	var b strings.Builder
	b.WriteString("(if ")
	b.WriteString(printer.expressionPrinter.Print(stmt.Condition))

	b.WriteString(" ")
	b.WriteString(printer.Print(stmt.ThenBranch))
	if stmt.ElseBranch != nil {
		b.WriteString(" ")
		b.WriteString(printer.Print(stmt.ElseBranch))
	}
	b.WriteString(")")
	return b.String()
}

func (printer *StatementPrinter) VisitWhileStatement(stmt *WhileStatement) any {
	var b strings.Builder
	b.WriteString("(while ")
	b.WriteString(printer.expressionPrinter.Print(stmt.Condition))

	b.WriteString(" ")
	b.WriteString(printer.Print(stmt.Body))
	b.WriteString(")")
	return b.String()
}

func (printer *StatementPrinter) VisitFunctionStatement(stmt *FunctionStatement) any {
	var b strings.Builder
	b.WriteString("(define (")
	b.WriteString(stmt.Name.Lexeme)
	for _, param := range stmt.Parameters {
		b.WriteString(" ")
		b.WriteString(param.Lexeme)
	}
	b.WriteString(")\n")

	for _, s := range stmt.Body.Statements {
		b.WriteString(printer.Print(s))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}

func (printer *StatementPrinter) VisitReturnStatement(stmt *ReturnStatement) any {
	return fmt.Sprintf("(return %s)", stmt.Value.Accept(printer.expressionPrinter))
}
