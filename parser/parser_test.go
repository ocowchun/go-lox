package parser

import (
	"testing"

	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/lexer"
)

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		// {"empty", "", ""},
		{"number", "1;", "1"},
		{"plus expression", "1 + 2;", "(+ 1 2)"},
		{"print statement", "print 1 + 2;", "(print (+ 1 2))"},
		{"var statement", "var a = 123;", "(define a 123)"},
		{"block statement", "{ var a = 123; print a;}", "(begin\n(define a 123)\n(print a)\n)"},
		{"if statement", "if (1 > 2) { print 1; }", "(if (> 1 2) (begin\n(print 1)\n))"},
		{"if else statement", "if (a > b) { print a; } else { print b; }", "(if (> a b) (begin\n(print a)\n) (begin\n(print b)\n))"},
		{"while statement", "while (i < 5) { i = i + 1;}", "(while (< i 5) (begin\n(set! i (+ i 1))\n))"},
		{"for statement", "for (var i = 0; i < 5; i = i + 1) { print i;}", "(begin\n(define i 0)\n(while (< i 5) (begin\n(begin\n(print i)\n)\n(set! i (+ i 1))\n))\n)"},
		{"function statement", "fun foo(a, b) { print a + b; }", "(define (foo a b)\n(print (+ a b))\n)"},
		{"return statement", "return 1 + 2;", "(return (+ 1 2))"},
		{"class statement", "class Foo { bar() { print 123; } }", "(class Foo\n(define (bar)\n(print 123)\n)\n)"},
		{"class statement with super class", "class Foo < Bar { bar() { print 123; } }", "(class Foo < Bar\n(define (bar)\n(print 123)\n)\n)"},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			lex := lexer.New(testCase.input)
			tokens, err := lex.Tokens()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			p := NewParser(tokens)

			statements, err := p.Parse()
			if err != nil {
				t.Fatalf("Failed to parse %s, error: %v", testCase.input, err)
			}

			if len(statements) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(statements))
			}
			printer := ast.Printer{}
			actual := printer.PrintStatement(statements[0])
			if actual != testCase.expected {
				t.Errorf("Expected %s, got %s", testCase.expected, actual)
			}
		})
	}

}

func TestParser_parseExpression(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"number", "1", "1"},
		{"string", "\"hello\"", "hello"},
		{"variable", "foo", "foo"},
		{"equal expression", "1 == 2", "(== 1 2)"},
		{"not equal expression", "1 != 2", "(!= 1 2)"},
		{"greater than expression", "1 > 2", "(> 1 2)"},
		{"greater than or equal expression", "1 >= 2", "(>= 1 2)"},
		{"less than expression", "1 < 2", "(< 1 2)"},
		{"less than or equal expression", "1 <= 2", "(<= 1 2)"},
		{"plus expression", "1 + 2", "(+ 1 2)"},
		{"minus expression", "1 - 2", "(- 1 2)"},
		{"multiply expression", "1 * 2", "(* 1 2)"},
		{"divide expression", "1 / 2", "(/ 1 2)"},
		{"bang expression", "!true", "(! true)"},
		{"negative expression", "-1", "(- 1)"},
		{"grouping expression", "(1 + 2)", "(group (+ 1 2))"},
		{"different precedence case 1", "1 + 2 * 3 - 4", "(- (+ 1 (* 2 3)) 4)"},
		{"different precedence case 2", "1 > 2 != 2 > 3", "(!= (> 1 2) (> 2 3))"},
		{"comma operator", "1 + 1, 2", "(begin (+ 1 1) 2)"},
		{"ternary operator", "1 > 2 ? 1 : 2", "(if (> 1 2) 1 2)"},
		{"assignment expression", "x = 1 + 2", "(set! x (+ 1 2))"},
		{"or expression", "a == b or a == c", "(or (== a b) (== a c))"},
		{"and expression", "a == b and a == c", "(and (== a b) (== a c))"},
		{"call expression 0", "foo()", "(foo)"},
		{"call expression 1", "foo(1)", "(foo 1)"},
		{"call expression 2", "foo(1, 2)", "(foo 1 2)"},
		{"function expression", "fun (a) { print a; }", "(lambda (a) (begin\n(print a)\n))"},
		{"get expression", "a.b", "(get a b)"},
		{"this expression", "this", "(this)"},
		{"super expression", "super.foo", "(super foo)"},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			lex := lexer.New(testCase.input)
			tokens, err := lex.Tokens()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			p := NewParser(tokens)

			expr, err := p.parseExpression()
			if err != nil {
				t.Fatalf("Failed to parse %s, error: %v", testCase.input, err)
			}

			printer := ast.Printer{}
			actual := printer.PrintExpression(expr)
			if actual != testCase.expected {
				t.Errorf("Expected %s, got %s", testCase.expected, actual)
			}
		})
	}
}

func TestParseInvalidExpression(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"number", "1 + !", "1"},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			lex := lexer.New(testCase.input)
			tokens, err := lex.Tokens()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			p := NewParser(tokens)

			_, err = p.Parse()
			if err == nil {
				t.Fatalf("Expected error for input %s, but got none", testCase.input)
			}
		})
	}
}
