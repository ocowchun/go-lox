package ast

import (
	"testing"

	"github.com/ocowchun/go-lox/token"
)

func TestBinaryExpression(t *testing.T) {
	lit1 := LiteralExpression{
		Value: "hello",
	}
	lit2 := LiteralExpression{
		Value: "world",
	}
	exp := BinaryExpression{
		Left:     &lit1,
		Operator: token.Token{Type: token.TokenTypePlus, Lexeme: "+"},
		Right:    &lit2,
	}
	printer := AstPrinter{}

	result := printer.Print(&exp)

	if result != "(+ hello world)" {
		t.Fatalf("Expected '(+ hello world)', got %v", result)
	}
}

func TestGroupedExpression(t *testing.T) {
	exp := GroupingExpression{
		Expression: &LiteralExpression{
			Value: "hello world",
		},
	}
	printer := AstPrinter{}

	result := printer.Print(&exp)

	if result != "(group hello world)" {
		t.Fatalf("Expected '(group hello world)', got %v", result)
	}

}

func TestLiteralExpression(t *testing.T) {
	exp := LiteralExpression{
		Value: "hello world",
	}
	printer := AstPrinter{}

	result := printer.Print(&exp)

	if result != "hello world" {
		t.Fatalf("Expected 'hello world', got %v", result)
	}
}

func TestUnaryExpression(t *testing.T) {
	exp := UnaryExpression{
		Operator: token.Token{Type: token.TokenTypeMinus, Literal: "-"},
		Right:    &LiteralExpression{Value: 123},
	}
	printer := AstPrinter{}

	result := printer.Print(&exp)

	if result != "(- 123)" {
		t.Fatalf("Expected '(- 123)', got %v", result)
	}

}
