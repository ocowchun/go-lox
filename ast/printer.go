package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

// Statement

func (printer *Printer) PrintStatement(stmt Stmt) string {
	res := stmt.Accept(printer).(string)
	return res
}

func (printer *Printer) VisitExpressionStatement(stmt *ExpressionStatement) any {
	return stmt.Expression.Accept(printer)
}

func (printer *Printer) VisitPrintStatement(stmt *PrintStatement) any {
	return fmt.Sprintf("(print %s)", stmt.Expression.Accept(printer))
}

func (printer *Printer) VisitVarStatement(stmt *VarStatement) any {
	return fmt.Sprintf("(define %s %s)", stmt.Name.Lexeme, stmt.Initializer.Accept(printer))
}

func (printer *Printer) VisitBlockStatement(stmt *BlockStatement) any {
	var b strings.Builder
	b.WriteString("(begin\n")
	for _, s := range stmt.Statements {
		b.WriteString(printer.PrintStatement(s))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}

func (printer *Printer) VisitIfStatement(stmt *IfStatement) any {
	var b strings.Builder
	b.WriteString("(if ")
	b.WriteString(printer.PrintExpression(stmt.Condition))

	b.WriteString(" ")
	b.WriteString(printer.PrintStatement(stmt.ThenBranch))
	if stmt.ElseBranch != nil {
		b.WriteString(" ")
		b.WriteString(printer.PrintStatement(stmt.ElseBranch))
	}
	b.WriteString(")")
	return b.String()
}

func (printer *Printer) VisitWhileStatement(stmt *WhileStatement) any {
	var b strings.Builder
	b.WriteString("(while ")
	b.WriteString(printer.PrintExpression(stmt.Condition))

	b.WriteString(" ")
	b.WriteString(printer.PrintStatement(stmt.Body))
	b.WriteString(")")
	return b.String()
}

func (printer *Printer) VisitFunctionStatement(stmt *FunctionStatement) any {
	var b strings.Builder
	b.WriteString("(define (")
	b.WriteString(stmt.Name.Lexeme)
	for _, param := range stmt.Parameters {
		b.WriteString(" ")
		b.WriteString(param.Lexeme)
	}
	b.WriteString(")\n")

	for _, s := range stmt.Body.Statements {
		b.WriteString(printer.PrintStatement(s))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}

func (printer *Printer) VisitReturnStatement(stmt *ReturnStatement) any {
	return fmt.Sprintf("(return %s)", stmt.Value.Accept(printer))
}

func (printer *Printer) VisitClassStatement(stmt *ClassStatement) any {
	// it's verbose to print class statements in a way that is similar to the Scheme syntax,
	var b strings.Builder
	b.WriteString("(class ")
	b.WriteString(stmt.Name.Lexeme)
	b.WriteString("\n")
	for _, method := range stmt.Methods {
		b.WriteString(printer.PrintStatement(method))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}

// Expression

func (printer *Printer) PrintExpression(expr Expr) string {
	res := expr.Accept(printer).(string)

	return res
}

func (printer *Printer) VisitBinaryExpression(expr *BinaryExpression) any {
	return fmt.Sprintf("(%s %s %s)",
		expr.Operator.Lexeme,
		printer.PrintExpression(expr.Left),
		printer.PrintExpression(expr.Right),
	)
}

func (printer *Printer) VisitGroupingExpression(expr *GroupingExpression) any {
	return fmt.Sprintf("(group %s)", printer.PrintExpression(expr.Expression))
}

func (printer *Printer) VisitLiteralExpression(expr *LiteralExpression) any {
	if str, ok := expr.Value.(string); ok {
		return str
	} else if num, ok := expr.Value.(float64); ok {
		return strconv.FormatFloat(num, 'f', -1, 64)
	} else {
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (printer *Printer) VisitUnaryExpression(expr *UnaryExpression) any {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, printer.PrintExpression(expr.Right))
}

func (printer *Printer) VisitCommaExpression(expr *CommaExpression) any {
	var b strings.Builder

	b.WriteString("(begin")
	for _, e := range expr.Expressions {
		b.WriteString(" ")
		b.WriteString(printer.PrintExpression(e))
	}
	b.WriteString(")")

	return b.String()
}

func (printer *Printer) VisitConditionExpression(expr *ConditionExpression) any {
	var b strings.Builder

	b.WriteString("(if ")
	b.WriteString(printer.PrintExpression(expr.Predicate))
	b.WriteString(" ")
	b.WriteString(printer.PrintExpression(expr.Consequent))
	b.WriteString(" ")
	b.WriteString(printer.PrintExpression(expr.Alternative))
	b.WriteString(")")

	return b.String()
}

func (printer *Printer) VisitVariableExpression(expr *VariableExpression) any {
	return expr.Name.Lexeme
}

func (printer *Printer) VisitAssignExpression(expr *AssignExpression) any {
	return fmt.Sprintf("(set! %s %s)", expr.Name.Lexeme, printer.PrintExpression(expr.Value))
}

func (printer *Printer) VisitLogicalExpression(expr *LogicalExpression) any {
	return fmt.Sprintf("(%s %s %s)",
		expr.Operator.Lexeme,
		printer.PrintExpression(expr.Left),
		printer.PrintExpression(expr.Right),
	)
}

func (printer *Printer) VisitCallExpression(expr *CallExpression) any {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(printer.PrintExpression(expr.Callee))

	for _, arg := range expr.Arguments {
		b.WriteString(" ")
		b.WriteString(printer.PrintExpression(arg))
	}
	b.WriteString(")")
	return b.String()
}

// (lambda (x y) (+ x y))
func (printer *Printer) VisitFunctionExpression(expr *FunctionExpression) any {
	var b strings.Builder
	b.WriteString("(lambda (")

	for i, parameter := range expr.Parameters {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(parameter.Lexeme)
	}
	b.WriteString(") ")
	b.WriteString(printer.PrintStatement(expr.Body))
	b.WriteString(")")
	return b.String()
}

func (printer *Printer) VisitGetExpression(expr *GetExpression) any {
	return fmt.Sprintf("(get %s %s)", printer.PrintExpression(expr.Object), expr.Name.Lexeme)
}

func (printer *Printer) VisitSetExpression(expr *SetExpression) any {
	return fmt.Sprintf("(set! %s %s %s)",
		printer.PrintExpression(expr.Object),
		expr.Name.Lexeme,
		printer.PrintExpression(expr.Value),
	)
}

func (printer *Printer) VisitThisExpression(expr *ThisExpression) any {
	return "(this)"
}
