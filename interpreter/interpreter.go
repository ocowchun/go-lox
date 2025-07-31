package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/token"
)

type Interpreter struct {
}

func New() *Interpreter {
	return &Interpreter{}
}

type EvaluatedResult struct {
	Value any
	Error error
}

func (interpreter *Interpreter) Evaluate(expr ast.Expr) EvaluatedResult {
	res := expr.Accept(interpreter).(EvaluatedResult)

	return res
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func NewRuntimeError(token token.Token, message string) *RuntimeError {
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (e *RuntimeError) Error() string {
	return e.Message
}

func (interpreter *Interpreter) VisitBinaryExpression(expr *ast.BinaryExpression) any {
	left := interpreter.Evaluate(expr.Left)
	if left.Error != nil {
		return EvaluatedResult{Error: left.Error}
	}

	right := interpreter.Evaluate(expr.Right)
	if right.Error != nil {
		return EvaluatedResult{Error: right.Error}
	}

	switch expr.Operator.Type {
	case token.TokenTypePlus:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue + rightValue}
			}
		} else if leftValue, ok := left.Value.(string); ok {
			if rightValue, ok := right.Value.(string); ok {
				return EvaluatedResult{Value: leftValue + rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers/strings for addition, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeMinus:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue - rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for subtraction, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeSlash:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				if rightValue == 0 {
					runtimeErr := NewRuntimeError(
						expr.Operator,
						"division by zero is not allowed",
					)
					return EvaluatedResult{Error: runtimeErr}
				}
				return EvaluatedResult{Value: leftValue / rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for division, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeStar:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue * rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for multiplication, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeGreater:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue > rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for greater than comparison, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeGreaterEqual:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue >= rightValue}
			}
		}
		return EvaluatedResult{Error: fmt.Errorf("expected numbers for greater than or equal comparison, got %T and %T", left.Value, right.Value)}
	case token.TokenTypeLess:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue < rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for less than comparison, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeLessEqual:
		if leftValue, ok := left.Value.(float64); ok {
			if rightValue, ok := right.Value.(float64); ok {
				return EvaluatedResult{Value: leftValue <= rightValue}
			}
		}

		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("expected numbers for less than or equal comparison, got %T and %T", left.Value, right.Value),
		)
		return EvaluatedResult{Error: runtimeErr}

	case token.TokenTypeEqualEqual:
		return EvaluatedResult{Value: isEqual(left.Value, right.Value)}

	case token.TokenTypeBangEqual:
		return EvaluatedResult{Value: isEqual(left.Value, right.Value)}

	default:
		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("unknown binary operator: %s", expr.Operator.Lexeme),
		)
		return EvaluatedResult{Error: runtimeErr}
	}
}

func (interpreter *Interpreter) VisitGroupingExpression(expr *ast.GroupingExpression) any {
	return interpreter.Evaluate(expr.Expression)
}

func (interpreter *Interpreter) VisitLiteralExpression(expr *ast.LiteralExpression) any {
	return EvaluatedResult{Value: expr.Value}
}

func (interpreter *Interpreter) VisitUnaryExpression(expr *ast.UnaryExpression) any {
	right := interpreter.Evaluate(expr.Right)
	if right.Error != nil {
		return EvaluatedResult{Error: right.Error}
	}

	switch expr.Operator.Type {
	case token.TokenTypeMinus:
		if value, ok := right.Value.(float64); ok {
			return EvaluatedResult{Value: -value}
		} else {
			runtimeErr := NewRuntimeError(
				expr.Operator,
				fmt.Sprintf("expected a number for unary minus, got %T", right.Value),
			)
			return EvaluatedResult{Error: runtimeErr}
		}
	case token.TokenTypeBang:
		return EvaluatedResult{Value: !isTruthy(right.Value)}

	default:
		runtimeErr := NewRuntimeError(
			expr.Operator,
			fmt.Sprintf("unknown unary operator: %s", expr.Operator.Lexeme),
		)
		return EvaluatedResult{Error: runtimeErr}
	}
}

func isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	if leftFloat, ok := left.(float64); ok {
		if rightFloat, ok := right.(float64); ok {
			return leftFloat == rightFloat
		}
	}

	if leftString, ok := left.(string); ok {
		if rightString, ok := right.(string); ok {
			return leftString == rightString
		}
	}
	if leftBool, ok := left.(bool); ok {
		if rightBool, ok := right.(bool); ok {
			return leftBool == rightBool
		}
	}

	return false
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}

	if boolean, ok := val.(bool); ok {
		return boolean
	}

	return true
}

func (interpreter *Interpreter) VisitBeginExpression(expr *ast.BeginExpression) any {
	// TODO
	return nil
	//var b strings.Builder
	//
	//b.WriteString("(begin ")
	//b.WriteString(interpreter.Evaluate(expr.Left))
	//b.WriteString(" ")
	//b.WriteString(interpreter.Evaluate(expr.Right))
	//b.WriteString(")")
	//
	//return b.String()
}

func (interpreter *Interpreter) VisitConditionExpression(expr *ast.ConditionExpression) any {
	// TODO
	return nil
	//var b strings.Builder
	//
	//b.WriteString("(if ")
	//b.WriteString(interpreter.Evaluate(expr.Predicate))
	//b.WriteString(" ")
	//b.WriteString(interpreter.Evaluate(expr.Consequent))
	//b.WriteString(" ")
	//b.WriteString(interpreter.Evaluate(expr.Alternative))
	//b.WriteString(")")
	//
	//return b.String()
}
