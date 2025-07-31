package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ocowchun/go-lox/interpreter"
	"github.com/ocowchun/go-lox/parser"
	"io"
	"os"
	"strings"

	"github.com/ocowchun/go-lox/lexer"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		target := args[1]
		runFile(target)

	} else if len(args) == 1 {
		runPrompt()

	} else {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	}
}

func runFile(target string) {
	file, err := os.Open(target)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(65)
	}
	defer file.Close()

	val, err := run(file)

	if err != nil {
		var runtimeError *interpreter.RuntimeError
		if errors.As(err, &runtimeError) {
			fmt.Printf("%s\n[line %d]\n", runtimeError.Message, runtimeError.Token.Line)
			os.Exit(70)
		} else {
			fmt.Println(err)
			os.Exit(65)
		}
	} else {
		if val != nil {
			fmt.Println(val)
		} else {
			fmt.Println("nil")
		}
	}
	// fmt.Println("Running file:", target)
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Running REPL")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}
		val, err := run(strings.NewReader(line))
		if err != nil {
			var runtimeError *interpreter.RuntimeError
			if errors.As(err, &runtimeError) {
				fmt.Printf("%s\n[line %d]\n", runtimeError.Message, runtimeError.Token.Line)
			} else {
				fmt.Println(err)
			}
		} else {
			if val != nil {
				fmt.Println(val)
			} else {
				fmt.Println("nil")
			}
		}
	}
	fmt.Println("Goodbye!")
}

func run(r io.Reader) (any, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	lex := lexer.New(buf.String())

	tokens, err := lex.Tokens()
	if err != nil {
		return nil, fmt.Errorf("lexer error: %s", err)
	}
	p := parser.NewParser(tokens)

	expr, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}

	i := interpreter.New()
	res := i.Evaluate(expr)
	return res.Value, res.Error
}
