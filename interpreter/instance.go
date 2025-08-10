package interpreter

import (
	"fmt"
	"github.com/ocowchun/go-lox/token"
)

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}

func (i *Instance) Get(name token.Token) (any, error) {
	if value, exists := i.fields[name.Lexeme]; exists {
		return value, nil
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i), nil
	}

	return nil, fmt.Errorf("undefined property '%s' in instance of class '%s'", name.Lexeme, i.class.name)
}

func (i *Instance) Set(name token.Token, value any) {
	i.fields[name.Lexeme] = value
}
