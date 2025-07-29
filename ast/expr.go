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

type BeginExpression struct {
	//Expressions []Expr
	Left  Expr
	Right Expr
}

func (exp *BeginExpression) Expr() {}

func (exp *BeginExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitBeginExpression(exp)
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

type ExprVisitor interface {
	VisitBinaryExpression(expr *BinaryExpression) any
	VisitGroupingExpression(expr *GroupingExpression) any
	VisitLiteralExpression(expr *LiteralExpression) any
	VisitUnaryExpression(expr *UnaryExpression) any
	VisitBeginExpression(expr *BeginExpression) any
	VisitConditionExpression(expr *ConditionExpression) any
}

type AstPrinter struct {
}

func (printer *AstPrinter) Print(expr Expr) string {
	res := expr.Accept(printer).(string)

	return res
}

func (printer *AstPrinter) VisitBinaryExpression(expr *BinaryExpression) any {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, printer.Print(expr.Left), printer.Print(expr.Right))

}

func (printer *AstPrinter) VisitGroupingExpression(expr *GroupingExpression) any {
	return fmt.Sprintf("(group %s)", printer.Print(expr.Expression))
}

func (printer *AstPrinter) VisitLiteralExpression(expr *LiteralExpression) any {
	if str, ok := expr.Value.(string); ok {
		return str
	} else if num, ok := expr.Value.(float64); ok {
		return strconv.FormatFloat(num, 'f', -1, 64)
	} else {
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (printer *AstPrinter) VisitUnaryExpression(expr *UnaryExpression) any {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, printer.Print(expr.Right))
}

func (printer *AstPrinter) VisitBeginExpression(expr *BeginExpression) any {
	var b strings.Builder

	b.WriteString("(begin ")
	b.WriteString(printer.Print(expr.Left))
	b.WriteString(" ")
	b.WriteString(printer.Print(expr.Right))
	b.WriteString(")")

	return b.String()
}

func (printer *AstPrinter) VisitConditionExpression(expr *ConditionExpression) any {
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
