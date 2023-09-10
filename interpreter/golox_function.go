package interpreter

import (
	"fmt"
	"golox/statement"
)

type GoloxFunction struct {
	Declaration statement.Function
}

func (g *GoloxFunction) Call(
	fInterpreter *Interpreter,
	arguments []any,
) any {
	environment := NewEnvironment(
		fInterpreter.Globals,
	)

	for i := 0; i < len(g.Declaration.Params); i++ {
		environment.Define(
			g.Declaration.Params[i].Lexeme,
			arguments[i],
		)
	}

	var res any = nil

	res = fInterpreter.ExecuteBlock(
		g.Declaration.Body,
		environment,
	)

	return res
}

func (g *GoloxFunction) Arity() int {
	return len(g.Declaration.Params)
}

func (g *GoloxFunction) ToString() string {
	return fmt.Sprintf("<fn %v>", g.Declaration.Name.Lexeme)
}
