This is an interpreter project for the language Lox from [craftinginterpreters](craftinginterpreters.com)
built in Golang, hence the name Golox. To run this project, you need to have Golang installed on your machine.

The goal is to make a working interpreter. Currently, the interpreter consists of :

- Scanner
- Parser
- Interpreter

The interpreter is currently able to evaluate simple mathematical expressions. 

Hit up the REPL (Read-Eval-Print-Loop) by using the command `go run main.go`. The interpreter supports file inputs by passing the file name to the run command `go run main.go [filename]` but is currently only able to evaluate a single expression. 