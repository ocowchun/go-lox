package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name token.Token, value any) error {
	if _, exists := e.values[name.Lexeme]; !exists {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
	
		return NewRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
	}

	e.values[name.Lexeme] = value
	return nil
}

func (e *Environment) Get(name token.Token) (any, error) {
	value, exists := e.values[name.Lexeme]
	if !exists {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}

		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
	}
	return value, nil
}
