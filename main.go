package main

import (
	"bufio"
	"fmt"
	"golox/ast"
	"golox/interpreter"
	"golox/parser"
	"golox/scanner"
	"os"
)

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

func printAst(expr ast.Expr) {
	if expr == nil {
		return
	}

	if v, ok := expr.(*ast.Literal); ok {
		fmt.Println(v)
	}

	if v, ok := expr.(*ast.Binary); ok {
		printAst(v.Left)
		printAst(v.Right)
		fmt.Println(v.Operator)
	}
}

func run(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()

	parser := parser.Parser{
		Tokens: tokens,
	}
	expression := parser.Parse()

	interpreter := interpreter.Interpreter{}

	interpreter.Interpret(expression)

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }
	// var expr ast.Expr = expression
	// for expression != nil {
	// 	fmt.Println(expr)
	// 	if v, ok := expr.(*ast.Binary); ok {
	// 		fmt.Println("left : ", (v.Left))
	// 		fmt.Println("right : ", (v.Right))
	// 		expr = v.Right
	// 	} else if v, ok := expr.(*ast.Unary); ok {
	// 		expr = v.Right
	// 	} else {
	// 		fmt.Println(expr)
	// 		break
	// 	}
	// }

	// printAst(expression)
}
