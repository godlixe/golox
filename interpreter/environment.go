package interpreter

import (
	"fmt"
	"golox/token"
)

type Environment struct {

	// Reference for nested environment(s).
	Enclosing *Environment

	// Values contains the variables in the golox
	// environment. Beware that the map is nil
	// and needs to be initialized before using.
	Values map[string]any
}

func NewEnvironment(enclosing Environment) Environment {
	return Environment{
		Enclosing: &enclosing,
		Values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name token.Token) any {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	// TODO return error
	fmt.Println("Undefined variable '", name.Lexeme, "'.")
	return nil
}

func (e *Environment) assign(name token.Token, value any) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.assign(name, value)
		return
	}

	fmt.Println("Undefined variable '", name.Lexeme, "'.")
}
