package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/token"
	"time"
)

type Interpreter struct {
	environment *Environment
}

// TODO: move builtin to a separate file
type clockFunction struct {
}

func (c *clockFunction) Call(interpreter *Interpreter, args []any) EvaluatedResult {
	return EvaluatedResult{
		Value: float64(time.Now().Unix()),
	}
}

func (c *clockFunction) Arity() int {
	return 0
}

func New() *Interpreter {
	global := NewEnvironment(nil)

	global.Define("clock", &clockFunction{})

	return &Interpreter{
		environment: global,
	}
}

type EvaluatedResult struct {
	Value any
	Error error
}

func (interpreter *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, stmt := range statements {
		res := interpreter.execute(stmt)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

type StatementResult struct {
	Value any
	Error error
}

func (interpreter *Interpreter) execute(statement ast.Stmt) StatementResult {
	res := statement.Accept(interpreter).(StatementResult)
	return res
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

func (interpreter *Interpreter) VisitWhileStatement(stmt *ast.WhileStatement) any {
	for {
		cond := interpreter.Evaluate(stmt.Condition)
		if cond.Error != nil {
			return cond.Error
		}

		if !isTruthy(cond.Value) {
			break
		}

		res := interpreter.execute(stmt.Body)
		if res.Error != nil {
			return res
		}
	}

	return StatementResult{}
}

func (interpreter *Interpreter) VisitIfStatement(stmt *ast.IfStatement) any {
	cond := interpreter.Evaluate(stmt.Condition)
	if cond.Error != nil {
		return cond.Error
	}

	if isTruthy(cond.Value) {
		return interpreter.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return interpreter.execute(stmt.ElseBranch)
	}

	return StatementResult{}
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

	return StatementResult{}
}

func (interpreter *Interpreter) VisitBlockStatement(stmt *ast.BlockStatement) any {
	res := interpreter.executeBlockStatement(stmt, NewEnvironment(interpreter.environment))

	return res
}

func (interpreter *Interpreter) executeBlockStatement(stmt *ast.BlockStatement, environment *Environment) StatementResult {
	// TODO: change to pass environment as a parameter to all visit methods
	previousEnvironment := interpreter.environment
	interpreter.environment = environment

	defer func() {
		interpreter.environment = previousEnvironment
	}()

	for _, statement := range stmt.Statements {
		res := interpreter.execute(statement)
		if res.Error != nil {
			return res
		} else if _, ok := res.Value.(ReturnValue); ok {
			return res
		}
	}

	return StatementResult{}
}

type Function struct {
	declaration *ast.FunctionStatement
	closure     *Environment // The environment in which the function was defined
}

func NewFunction(declaration *ast.FunctionStatement, closure *Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *Function) Call(interpreter *Interpreter, args []any) EvaluatedResult {
	environment := NewEnvironment(f.closure)

	if len(args) != f.Arity() {
		return EvaluatedResult{
			Error: NewRuntimeError(
				f.declaration.Name,
				fmt.Sprintf("expected %d arguments but got %d", f.Arity(), len(args)),
			),
		}
	}

	for i, param := range f.declaration.Parameters {
		environment.Define(param.Lexeme, args[i])
	}

	// TODO: handle return value
	res := interpreter.executeBlockStatement(f.declaration.Body, environment)
	if res.Error != nil {
		return EvaluatedResult{Error: res.Error}
	}

	if returnValue, ok := res.Value.(ReturnValue); ok {
		return EvaluatedResult{
			Value: returnValue.Value,
		}

	} else {
		// If no return value is specified, return nil
		return EvaluatedResult{
			Value: nil,
		}
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.Parameters)
}

func (f *Function) String() string {
	printer := ast.NewStatementPrinter()
	return printer.Print(f.declaration)
}

func (interpreter *Interpreter) VisitFunctionStatement(stmt *ast.FunctionStatement) any {
	function := NewFunction(stmt, interpreter.environment)
	interpreter.environment.Define(stmt.Name.Lexeme, function)

	return StatementResult{
		Error: nil,
	}
}

func (interpreter *Interpreter) VisitExpressionStatement(stmt *ast.ExpressionStatement) any {
	result := interpreter.Evaluate(stmt.Expression)
	return StatementResult{
		Error: result.Error,
	}
}

type ReturnValue struct {
	Value any
}

func (interpreter *Interpreter) VisitReturnStatement(stmt *ast.ReturnStatement) any {
	result := interpreter.Evaluate(stmt.Value)

	return StatementResult{
		Value: ReturnValue{Value: result.Value},
		Error: result.Error,
	}
}

func (interpreter *Interpreter) VisitPrintStatement(stmt *ast.PrintStatement) any {
	result := interpreter.Evaluate(stmt.Expression)
	if result.Error != nil {
		return StatementResult{Error: result.Error}
	}

	if result.Value != nil {
		fmt.Println(result.Value)
	} else {
		fmt.Println("nil")
	}

	return StatementResult{}
}

func (interpreter *Interpreter) VisitLogicalExpression(expr *ast.LogicalExpression) any {
	left := interpreter.Evaluate(expr.Left)
	if left.Error != nil {
		return left
	}

	if expr.Operator.Type == token.TokenTypeOr {
		if isTruthy(left.Value) {
			return left
		}
	} else {
		if !isTruthy(left.Value) {
			return left
		}
	}

	return interpreter.Evaluate(expr.Right)
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

func (interpreter *Interpreter) VisitCommaExpression(expr *ast.CommaExpression) any {
	fmt.Println("===VisitCommaExpression", expr)
	var res EvaluatedResult
	for _, subExpr := range expr.Expressions {
		result := interpreter.Evaluate(subExpr)
		if result.Error != nil {
			return result
		}

		// Update res with the last evaluated value
		res = result
	}

	return res
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

func (interpreter *Interpreter) VisitCallExpression(expr *ast.CallExpression) any {
	evaluatedResult := interpreter.Evaluate(expr.Callee)
	if evaluatedResult.Error != nil {
		return evaluatedResult
	}

	var function Callable
	if callable, ok := evaluatedResult.Value.(Callable); ok {
		function = callable
	} else {
		runtimeErr := NewRuntimeError(
			expr.Paren,
			fmt.Sprintf("can only call functions and classes, got %T", evaluatedResult.Value),
		)
		return EvaluatedResult{Error: runtimeErr}
	}

	if len(expr.Arguments) != function.Arity() {
		runtimeErr := NewRuntimeError(
			expr.Paren,
			fmt.Sprintf("expected %d arguments but got %d", function.Arity(), len(expr.Arguments)),
		)
		return EvaluatedResult{Error: runtimeErr}
	}

	args := make([]any, 0, len(expr.Arguments))
	for _, argExp := range expr.Arguments {
		evaluatedResult = interpreter.Evaluate(argExp)
		if evaluatedResult.Error != nil {
			return evaluatedResult
		}
		args = append(args, evaluatedResult.Value)
	}

	return function.Call(interpreter, args)
}

type Callable interface {
	Call(interpreter *Interpreter, args []any) EvaluatedResult
	Arity() int
}
