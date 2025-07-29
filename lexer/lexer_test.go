package lexer

import (
	"math"
	"testing"

	"github.com/ocowchun/go-lox/token"
)

func TestLexer(t *testing.T) {
	input := "( ) { } , . - + * ; ! != = == < <= > >= / 123 \"hello lexer\" foo and class else false for fun if nil or print return super this true var while"
	l := New(input)

	expectedTokens := []token.Token{
		token.Token{Type: token.TokenTypeLeftParen},
		token.Token{Type: token.TokenTypeRightParen},
		token.Token{Type: token.TokenTypeLeftBrace},
		token.Token{Type: token.TokenTypeRightBrace},
		token.Token{Type: token.TokenTypeComma},
		token.Token{Type: token.TokenTypeDot},
		token.Token{Type: token.TokenTypeMinus},
		token.Token{Type: token.TokenTypePlus},
		token.Token{Type: token.TokenTypeStar},
		token.Token{Type: token.TokenTypeSemicolon},
		token.Token{Type: token.TokenTypeBang},
		token.Token{Type: token.TokenTypeBangEqual},
		token.Token{Type: token.TokenTypeEqual},
		token.Token{Type: token.TokenTypeEqualEqual},
		token.Token{Type: token.TokenTypeLess},
		token.Token{Type: token.TokenTypeLessEqual},
		token.Token{Type: token.TokenTypeGreater},
		token.Token{Type: token.TokenTypeGreaterEqual},
		token.Token{Type: token.TokenTypeSlash},
		token.Token{Type: token.TokenTypeNumber, Literal: float64(123)},
		token.Token{Type: token.TokenTypeString, Literal: "hello lexer"},
		token.Token{Type: token.TokenTypeIdentifier, Lexeme: "foo"},
		token.Token{Type: token.TokenTypeAnd},
		token.Token{Type: token.TokenTypeClass},
		token.Token{Type: token.TokenTypeElse},
		token.Token{Type: token.TokenTypeFalse},
		token.Token{Type: token.TokenTypeFor},
		token.Token{Type: token.TokenTypeFun},
		token.Token{Type: token.TokenTypeIf},
		token.Token{Type: token.TokenTypeNil},
		token.Token{Type: token.TokenTypeOr},
		token.Token{Type: token.TokenTypePrint},
		token.Token{Type: token.TokenTypeReturn},
		token.Token{Type: token.TokenTypeSuper},
		token.Token{Type: token.TokenTypeThis},
		token.Token{Type: token.TokenTypeTrue},
		token.Token{Type: token.TokenTypeVar},
		token.Token{Type: token.TokenTypeWhile},
	}

	i := 0
	for !l.IsAtEnd() {
		token, err := l.NextToken()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		assertToken(t, token, expectedTokens[i])
		i++
	}

}

const float64EqualityThreshold = 1e-9

func assertToken(t *testing.T, actualToken token.Token, expectedToken token.Token) {
	if actualToken.Type != expectedToken.Type {
		t.Fatalf("Expected token to be %s, got %s", expectedToken.Type, actualToken.Type)
	}

	if actualToken.Type == token.TokenTypeNumber {
		actualLiteral, ok := actualToken.Literal.(float64)
		if !ok {
			t.Fatalf("failed to convert actual literal to float64, literal: %v", actualToken.Literal)
		}

		expectedLiteral, ok := expectedToken.Literal.(float64)
		if !ok {
			t.Fatalf("failed to convert expected literal to float64, literal: %v", expectedToken.Literal)
		}

		if math.Abs(actualLiteral-expectedLiteral) > float64EqualityThreshold {
			t.Fatalf("literal not matched, expected = %v, actual = %v  ", expectedLiteral, actualLiteral)

		}
	} else if actualToken.Type == token.TokenTypeString {
		actualLiteral, ok := actualToken.Literal.(string)
		if !ok {
			t.Fatalf("failed to convert actual literal to string, literal: %v", actualToken.Literal)
		}

		expectedLiteral, ok := expectedToken.Literal.(string)
		if !ok {
			t.Fatalf("failed to convert expected literal to string, literal: %v", expectedToken.Literal)
		}

		if actualLiteral != expectedLiteral {
			t.Fatalf("literal not matched, expected = %v, actual = %v  ", expectedLiteral, actualLiteral)
		}
	} else if actualToken.Type == token.TokenTypeIdentifier {
		if actualToken.Lexeme != expectedToken.Lexeme {
			t.Fatalf("Lexeme not matched, expected = %v, actual = %v  ", expectedToken.Lexeme, actualToken.Lexeme)
		}
	}
}
