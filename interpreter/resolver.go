package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/token"
)

type FunctionType uint8

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

type NameMetadata struct {
	// Whether the name is initialized
	initialized bool

	// Whether the name is used in the current/inner scope
	used bool
}

type Resolver struct {
	interpreter         *Interpreter
	scopes              []map[string]*NameMetadata
	currentFunctionType FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:         interpreter,
		scopes:              []map[string]*NameMetadata{},
		currentFunctionType: FunctionTypeNone,
	}
}

func (r *Resolver) ResolveStatement(statement ast.Stmt) error {
	err := statement.Accept(r)
	if err != nil {
		return err.(error)
	}
	return nil
}

func (r *Resolver) ResolveExpression(expr ast.Expr) error {
	err := expr.Accept(r)
	if err != nil {
		return err.(error)
	}
	return nil
}

func (r *Resolver) beginScope() {
	scope := make(map[string]*NameMetadata)
	r.scopes = append(r.scopes, scope)
}

func (r *Resolver) endScope() {
	if len(r.scopes) == 0 {
		panic("No scope to end")
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return NewResolveError(name, fmt.Sprintf("Already a variable with this name `%s` in this scope.", name.Lexeme))
	}
	scope[name.Lexeme] = &NameMetadata{
		initialized: false, // Mark as declared but not initialized
		used:        false, // Not used yet
	}

	return nil
}

func (r *Resolver) define(name token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme].initialized = true // Mark as initialized

	return nil
}

func (r *Resolver) VisitExpressionStatement(stmt *ast.ExpressionStatement) any {
	return r.ResolveExpression(stmt.Expression)
}

func (r *Resolver) VisitPrintStatement(stmt *ast.PrintStatement) any {
	return r.ResolveExpression(stmt.Expression)
}

func (r *Resolver) VisitVarStatement(stmt *ast.VarStatement) any {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}

	if stmt.Initializer != nil {
		err = r.ResolveExpression(stmt.Initializer)
		if err != nil {
			return err
		}
	}

	err = r.define(stmt.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitBlockStatement(stmt *ast.BlockStatement) any {
	r.beginScope()
	defer r.endScope()
	for _, s := range stmt.Statements {
		err := r.ResolveStatement(s)
		if err != nil {
			return err
		}
	}

	if r.currentFunctionType == FunctionTypeFunction {
		parametersScope := r.scopes[len(r.scopes)-2]
		blockScope := r.scopes[len(r.scopes)-1]
		for name, metadata := range blockScope {
			if _, ok := parametersScope[name]; ok {
				return NewResolveError(token.Token{Lexeme: name}, fmt.Sprintf("Local variable `%s` conflicts with parameter.", name))
			}

			if !metadata.used {
				return NewResolveError(token.Token{Lexeme: name}, fmt.Sprintf("Local variable `%s` is declared but never used.", name))
			}
		}

	}

	return nil
}

func (r *Resolver) VisitIfStatement(stmt *ast.IfStatement) any {
	err := r.ResolveExpression(stmt.Condition)
	if err != nil {
		return err
	}

	err = r.ResolveStatement(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		err = r.ResolveStatement(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitWhileStatement(stmt *ast.WhileStatement) any {
	err := r.ResolveExpression(stmt.Condition)
	if err != nil {
		return err
	}

	return r.ResolveStatement(stmt.Body)
}

func (r *Resolver) VisitFunctionStatement(stmt *ast.FunctionStatement) any {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}

	err = r.define(stmt.Name)
	if err != nil {
		return err
	}

	return r.resolveFunction(stmt.Parameters, stmt.Body, FunctionTypeFunction)
}

func (r *Resolver) resolveFunction(parameters []token.Token, body *ast.BlockStatement, functionType FunctionType) error {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = functionType

	r.beginScope()
	defer func() {
		r.currentFunctionType = enclosingFunctionType
		r.endScope()
	}()

	//parameter
	for _, param := range parameters {
		err := r.declare(param)
		if err != nil {
			return err
		}
		err = r.define(param)
		if err != nil {
			return err
		}
	}

	return r.ResolveStatement(body)
}

func (r *Resolver) VisitReturnStatement(stmt *ast.ReturnStatement) any {
	if r.currentFunctionType == FunctionTypeNone {
		return NewResolveError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		return r.ResolveExpression(stmt.Value)
	}

	return nil
}

// Expression

func (r *Resolver) VisitBinaryExpression(expr *ast.BinaryExpression) any {
	err := r.ResolveExpression(expr.Left)
	if err != nil {
		return err
	}

	err = r.ResolveExpression(expr.Right)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitGroupingExpression(expr *ast.GroupingExpression) any {
	return r.ResolveExpression(expr.Expression)
}

func (r *Resolver) VisitLiteralExpression(expr *ast.LiteralExpression) any {
	return nil
}

func (r *Resolver) VisitUnaryExpression(expr *ast.UnaryExpression) any {
	return r.ResolveExpression(expr.Right)
}

func (r *Resolver) VisitCommaExpression(expr *ast.CommaExpression) any {
	panic("TODO")
}

func (r *Resolver) VisitConditionExpression(expr *ast.ConditionExpression) any {
	panic("TODO")
}

type ResolveError struct {
	Token   token.Token
	Message string
}

func NewResolveError(token token.Token, message string) *ResolveError {
	return &ResolveError{
		Token:   token,
		Message: message,
	}
}

func (e *ResolveError) Error() string {
	return e.Message
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return nil
		}
	}
	return nil
}

func (r *Resolver) VisitVariableExpression(expr *ast.VariableExpression) any {
	if len(r.scopes) > 0 {
		metadata, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]
		if !ok {
			// Variable is not defined in the current scope
			// We assume it's a global variable
			return nil
		}
		if !metadata.initialized {
			return NewResolveError(expr.Name, "Can't read local variable in its own initializer.")
		}
		metadata.used = true
	}

	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) VisitAssignExpression(expr *ast.AssignExpression) any {
	// x = 1 + 2;
	err := r.ResolveExpression(expr.Value)
	if err != nil {
		return err
	}

	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) VisitLogicalExpression(expr *ast.LogicalExpression) any {
	err := r.ResolveExpression(expr.Left)
	if err != nil {
		return err
	}

	return r.ResolveExpression(expr.Right)
}

func (r *Resolver) VisitCallExpression(expr *ast.CallExpression) any {
	err := r.ResolveExpression(expr.Callee)
	if err != nil {
		return err
	}

	for _, arg := range expr.Arguments {
		err = r.ResolveExpression(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitFunctionExpression(expr *ast.FunctionExpression) any {
	return r.resolveFunction(expr.Parameters, expr.Body, FunctionTypeFunction)
}
