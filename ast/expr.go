package ast

import (
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

// LogicalExpression represents a logical operation, such as AND or OR.
// It is used to handle short-circuit evaluation in the interpreter.
// That's why we can't use BinaryExpression for this purpose, as it does not support short-circuiting.
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

type FunctionExpression struct {
	Fun        token.Token // keep the keyword for error reporting
	Parameters []token.Token
	Body       *BlockStatement
}

func (exp *FunctionExpression) Expr() {}

func (exp *FunctionExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitFunctionExpression(exp)
}

type GetExpression struct {
	Object Expr
	Name   token.Token
}

func (exp *GetExpression) Expr() {}

func (exp *GetExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitGetExpression(exp)
}

type SetExpression struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

func (exp *SetExpression) Expr() {}
func (exp *SetExpression) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetExpression(exp)
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
	VisitFunctionExpression(expr *FunctionExpression) any
	VisitGetExpression(expr *GetExpression) any
	VisitSetExpression(expr *SetExpression) any
}
