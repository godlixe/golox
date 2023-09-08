package interpreter

import (
	"fmt"
	"golox/statement"
)

type GoloxFunction struct {
	Declaration statement.Function
}

func (g *GoloxFunction) Call(
	gInterpreter *Interpreter,
	arguments []any,
) any {
	environment := NewEnvironment(
		gInterpreter.Globals,
	)

	for i := 0; i < len(g.Declaration.Params); i++ {
		environment.Define(
			g.Declaration.Params[i].Lexeme,
			arguments[i],
		)
	}

	gInterpreter.ExecuteBlock(
		g.Declaration.Body,
		environment,
	)

	return nil
}

func (g *GoloxFunction) Arity() int {
	return len(g.Declaration.Params)
}

func (g *GoloxFunction) ToString() string {
	return fmt.Sprintf("<fn %v>", g.Declaration.Name.Lexeme)
}
