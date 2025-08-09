package interpreter

import (
	"errors"
	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/lexer"
	"github.com/ocowchun/go-lox/parser"
	"testing"
)

func TestResolver_LocalVariableCannotBeDeclaredTwice(t *testing.T) {
	code := `
{
	var x = 1;
	var x = 2;
}`

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Already a variable with this name `x` in this scope." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_LocalVariableCannotShadowFunctionParameter(t *testing.T) {
	code := `
fun foo(x) {
	var x = 1;
}
`

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Local variable `x` conflicts with parameter." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_LocalVariableMustBeUsed(t *testing.T) {
	code := `
fun foo() {
    var a = 123;

}
`

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Local variable `a` is declared but never used." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_LocalVariableUsd(t *testing.T) {
	code := `
fun foo() {
    var a = 123;
	print a;
}
`

	err := resolveTestCode(code)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)

	}
}

func TestResolver_LocalVariableUsdInClosure(t *testing.T) {
	code := `
fun foo() {
    var a = 123;

    fun bar() {
        print a;
    }
    bar();
}
`

	err := resolveTestCode(code)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)

	}
}

func TestResolver_CannotReadFromOwnInitializer(t *testing.T) {
	// TODO: return error in top level code for var a = a;
	code := "{var a = a;}"

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Can't read local variable in its own initializer." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_CannotReturnFromTopLevel(t *testing.T) {
	code := `return 9527;`

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Can't return from top-level code." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_CannotReturnFromInitializer(t *testing.T) {
	code := `
class Foo {
	init() {
		return 123;
	}
}
`

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "Can't return a value from an initializer." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func TestResolver_ClassCannotInheritFromItself(t *testing.T) {
	code := "class Oops < Oops {}"

	err := resolveTestCode(code)

	var resolveError *ResolveError
	if !errors.As(err, &resolveError) {
		t.Fatalf("Expected ResolveError, got %T", err)
	} else {
		if resolveError.Message != "A class can't inherit from itself." {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

func resolveTestCode(code string) error {
	interpreter := New()
	resolver := NewResolver(interpreter)

	statements := parseCode(code)
	return resolver.ResolveStatements(statements)
}

func parseCode(code string) []ast.Stmt {
	l := lexer.New(code)
	tokens, err := l.Tokens()
	if err != nil {
		panic("Failed to tokenize code: " + err.Error())
	}

	p := parser.NewParser(tokens)
	statements, err := p.Parse()
	if err != nil {
		panic("Failed to parse code: " + err.Error())
	}
	return statements
}
