package main

import (
	"bufio"
	"fmt"
	"github.com/ocowchun/go-lox/ast"
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

	run(file)
	// fmt.Println("Running file:", target)
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}
		err := run(strings.NewReader(line))
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Goodbye!")
}

func run(r io.Reader) error {
	fmt.Println("Running REPL")
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	lex := lexer.New(buf.String())

	tokens, err := lex.Tokens()
	if err != nil {
		return fmt.Errorf("lexer error: %s", err)
	}
	p := parser.NewParser(tokens)

	expr, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %s", err)
	}

	printer := ast.AstPrinter{}
	fmt.Println(printer.Print(expr))

	return nil
}
