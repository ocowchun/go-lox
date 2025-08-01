package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/token"
)

type Interpreter struct {
	environment *Environment
}

func New() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

type EvaluatedResult struct {
	Value any
	Error error
}

func (interpreter *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, stmt := range statements {
		err := interpreter.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (interpreter *Interpreter) execute(statement ast.Stmt) error {
	err := statement.Accept(interpreter)
	if err != nil {
		return err.(error)
	}

	return nil
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

func (interpreter *Interpreter) VisitVarStatement(stmt *ast.VarStatement) any {
	if stmt.Initializer != nil {
		initResult := interpreter.Evaluate(stmt.Initializer)
		if initResult.Error != nil {
			return initResult.Error
		}
		interpreter.environment.Define(stmt.Name.Lexeme, initResult.Value)
	} else {
		interpreter.environment.Define(stmt.Name.Lexeme, nil)
	}

	return nil
}

func (interpreter *Interpreter) VisitBlockStatement(stmt *ast.BlockStatement) any {
	err := interpreter.executeBlockStatement(stmt, NewEnvironment(interpreter.environment))

	return err
}

func (interpreter *Interpreter) executeBlockStatement(stmt *ast.BlockStatement, environment *Environment) error {
	// TODO: change to pass environment as a parameter to all visit methods
	previousEnvironment := interpreter.environment
	interpreter.environment = environment

	defer func() {
		interpreter.environment = previousEnvironment
	}()

	for _, statement := range stmt.Statements {
		err := interpreter.execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (interpreter *Interpreter) VisitExpressionStatement(stmt *ast.ExpressionStatement) any {
	result := interpreter.Evaluate(stmt.Expression)
	return result.Error
}

func (interpreter *Interpreter) VisitPrintStatement(stmt *ast.PrintStatement) any {
	result := interpreter.Evaluate(stmt.Expression)
	if result.Error != nil {
		return result.Error
	}

	if result.Value != nil {
		fmt.Println(result.Value)
	} else {
		fmt.Println("nil")
	}

	return nil
}

func (interpreter *Interpreter) VisitVariableExpression(expr *ast.VariableExpression) any {
	val, err := interpreter.environment.Get(expr.Name)
	return EvaluatedResult{
		Value: val,
		Error: err,
	}
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

func (interpreter *Interpreter) VisitAssignExpression(expr *ast.AssignExpression) any {
	res := interpreter.Evaluate(expr.Value)
	if res.Error != nil {
		return res
	}

	err := interpreter.environment.Assign(expr.Name, res.Value)
	if err != nil {
		return EvaluatedResult{Error: err}
	}

	return res
}
