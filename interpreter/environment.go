package interpreter

import (
	"fmt"
	"golox/token"
)

type Environment struct {

	// Values contains the variables in the golox
	// environment. Beware that the map is nil
	// and needs to be initialized before using.
	Values map[string]any
}

func (e *Environment) define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) get(name token.Token) any {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v
	}

	// TODO return error
	fmt.Println("Undefined variable '", name.Lexeme)
	return nil
}
