package main

import (
	"bufio"
	"fmt"
	"golox/interpreter"
	"golox/parser"
	"golox/scanner"
	"os"
	"time"
)

// TODO : move to somewhere else
type clock struct{}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(
	interpreter interpreter.Interpreter,
	arguments []any,
) any {
	return float64(time.Now().UnixMilli() / 1000)
}

func (c *clock) ToString() string {
	return "<native fn>"
}

func main() {
	// get arguments from program
	args := os.Args

	// golox command expects 1 argument
	// which is the path of the script
	if len(args) > 2 {
		fmt.Println("Usage: golox [script]")
		return
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPromt()
	}
}

func runPromt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if text == "" {
			break
		}

		run(text)
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	run(string(data))
}

func run(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()

	parser := parser.Parser{
		Tokens: tokens,
	}
	statements := parser.Parse()

	// initialize global environment here for
	// a fixed reference to the outermost global
	// environment for the interpreter.
	globalEnv := interpreter.Environment{
		Enclosing: nil,
		Values:    make(map[string]any),
	}

	globalEnv.Define("clock", clock{})

	interpreter := interpreter.Interpreter{
		Environment: globalEnv,
		Globals:     globalEnv,
	}

	interpreter.Interpret(statements)
}
