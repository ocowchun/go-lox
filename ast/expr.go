package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ocowchun/go-lox/token"
)

type Expr interface {
	Expr()
	Accept(visitor ExprVisitor) any
}

type BinaryExpression struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (exp *BinaryExpression) Expr() {}

func (exp *BinaryExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpression(exp)
}

type GroupingExpression struct {
	Expression Expr
}

func (exp *GroupingExpression) Expr() {
}

func (exp *GroupingExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpression(exp)
}

type LiteralExpression struct {
	Value any
}

func (exp *LiteralExpression) Expr() {}

func (exp *LiteralExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpression(exp)
}

type UnaryExpression struct {
	Operator token.Token
	Right    Expr
}

func (exp *UnaryExpression) Expr() {
}

func (exp *UnaryExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpression(exp)
}

type CommaExpression struct {
	Expressions []Expr
}

func (exp *CommaExpression) Expr() {}

func (exp *CommaExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitCommaExpression(exp)
}

type ConditionExpression struct {
	Predicate   Expr
	Consequent  Expr
	Alternative Expr
}

func (exp *ConditionExpression) Expr() {}

func (exp *ConditionExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitConditionExpression(exp)
}

type VariableExpression struct {
	Name token.Token
}

func (exp *VariableExpression) Expr() {}

func (exp *VariableExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpression(exp)
}

type AssignExpression struct {
	Name  token.Token
	Value Expr
}

func (exp *AssignExpression) Expr() {}

func (exp *AssignExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpression(exp)
}

type LogicalExpression struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (exp *LogicalExpression) Expr() {}

func (exp *LogicalExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpression(exp)
}

type CallExpression struct {
	Callee Expr
	// For Runtime errors, we need to know the position of the opening parenthesis
	Paren     token.Token
	Arguments []Expr
}

func (exp *CallExpression) Expr() {}

func (exp *CallExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpression(exp)
}

type ExprVisitor interface {
	VisitBinaryExpression(expr *BinaryExpression) any
	VisitGroupingExpression(expr *GroupingExpression) any
	VisitLiteralExpression(expr *LiteralExpression) any
	VisitUnaryExpression(expr *UnaryExpression) any
	VisitCommaExpression(expr *CommaExpression) any
	VisitConditionExpression(expr *ConditionExpression) any
	VisitVariableExpression(expr *VariableExpression) any
	VisitAssignExpression(expr *AssignExpression) any
	VisitLogicalExpression(expr *LogicalExpression) any
	VisitCallExpression(expr *CallExpression) any
}

type ExpressionPrinter struct {
}

func (printer *ExpressionPrinter) Print(expr Expr) string {
	res := expr.Accept(printer).(string)

	return res
}

func (printer *ExpressionPrinter) VisitBinaryExpression(expr *BinaryExpression) any {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, printer.Print(expr.Left), printer.Print(expr.Right))

}

func (printer *ExpressionPrinter) VisitGroupingExpression(expr *GroupingExpression) any {
	return fmt.Sprintf("(group %s)", printer.Print(expr.Expression))
}

func (printer *ExpressionPrinter) VisitLiteralExpression(expr *LiteralExpression) any {
	if str, ok := expr.Value.(string); ok {
		return str
	} else if num, ok := expr.Value.(float64); ok {
		return strconv.FormatFloat(num, 'f', -1, 64)
	} else {
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (printer *ExpressionPrinter) VisitUnaryExpression(expr *UnaryExpression) any {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, printer.Print(expr.Right))
}

func (printer *ExpressionPrinter) VisitCommaExpression(expr *CommaExpression) any {
	var b strings.Builder

	b.WriteString("(begin")
	for _, e := range expr.Expressions {
		b.WriteString(" ")
		b.WriteString(printer.Print(e))
	}
	b.WriteString(")")

	return b.String()
}

func (printer *ExpressionPrinter) VisitConditionExpression(expr *ConditionExpression) any {
	var b strings.Builder

	b.WriteString("(if ")
	b.WriteString(printer.Print(expr.Predicate))
	b.WriteString(" ")
	b.WriteString(printer.Print(expr.Consequent))
	b.WriteString(" ")
	b.WriteString(printer.Print(expr.Alternative))
	b.WriteString(")")

	return b.String()
}

func (printer *ExpressionPrinter) VisitVariableExpression(expr *VariableExpression) any {
	return expr.Name.Lexeme
}

func (printer *ExpressionPrinter) VisitAssignExpression(expr *AssignExpression) any {
	return fmt.Sprintf("(set! %s %s)", expr.Name.Lexeme, printer.Print(expr.Value))
}

func (printer *ExpressionPrinter) VisitLogicalExpression(expr *LogicalExpression) any {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, printer.Print(expr.Left), printer.Print(expr.Right))
}

func (printer *ExpressionPrinter) VisitCallExpression(expr *CallExpression) any {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(printer.Print(expr.Callee))

	for _, arg := range expr.Arguments {
		b.WriteString(" ")
		b.WriteString(printer.Print(arg))
	}
	b.WriteString(")")
	return b.String()
}
