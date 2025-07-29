package lexer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ocowchun/go-lox/token"
)

type Lexer struct {
	source  string
	start   int
	current int
	line    int
}

func New(input string) *Lexer {
	return &Lexer{
		source:  input,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (l *Lexer) Tokens() ([]token.Token, error) {
	tokens := make([]token.Token, 0)

	for !l.IsAtEnd() {

		t, err := l.NextToken()
		if err != nil {
			return tokens, err
		}

		if t.IsTokenType(token.TokenTypeEOF) {
			break
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

func (l *Lexer) IsAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) Advance() byte {
	if l.IsAtEnd() {
		panic("can't called Advance when lexer is at end")
	}
	c := l.source[l.current]
	l.current++
	return c
}

func (l *Lexer) match(expected byte) bool {
	if l.IsAtEnd() {
		return false
	}

	if l.source[l.current] != expected {
		return false
	}

	l.current++
	return true
}

func (l *Lexer) peek() byte {
	if l.IsAtEnd() {
		return 0
	}

	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}

	return l.source[l.current+1]
}

func (l *Lexer) NextToken() (token.Token, error) {
	for !l.IsAtEnd() {
		l.start = l.current

		c := l.Advance()
		switch c {
		case '(':
			return token.Token{Type: token.TokenTypeLeftParen, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case ')':
			return token.Token{Type: token.TokenTypeRightParen, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '{':
			return token.Token{Type: token.TokenTypeLeftBrace, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '}':
			return token.Token{Type: token.TokenTypeRightBrace, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case ',':
			return token.Token{Type: token.TokenTypeComma, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '.':
			return token.Token{Type: token.TokenTypeDot, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '-':
			return token.Token{Type: token.TokenTypeMinus, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '+':
			return token.Token{Type: token.TokenTypePlus, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '*':
			return token.Token{Type: token.TokenTypeStar, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case ';':
			return token.Token{Type: token.TokenTypeSemicolon, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '?':
			return token.Token{Type: token.TokenTypeQuestionMark, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case ':':
			return token.Token{Type: token.TokenTypeColon, Lexeme: string(c), Literal: nil, Line: l.line}, nil
		case '!':
			if l.match('=') {
				return token.Token{Type: token.TokenTypeBangEqual, Lexeme: "!=", Literal: nil, Line: l.line}, nil
			} else {
				return token.Token{Type: token.TokenTypeBang, Lexeme: "!", Literal: nil, Line: l.line}, nil
			}
		case '=':
			if l.match('=') {
				return token.Token{Type: token.TokenTypeEqualEqual, Lexeme: "==", Literal: nil, Line: l.line}, nil
			} else {
				return token.Token{Type: token.TokenTypeEqual, Lexeme: "=", Literal: nil, Line: l.line}, nil
			}
		case '>':
			if l.match('=') {
				return token.Token{Type: token.TokenTypeGreaterEqual, Lexeme: ">=", Literal: nil, Line: l.line}, nil
			} else {
				return token.Token{Type: token.TokenTypeGreater, Lexeme: ">", Literal: nil, Line: l.line}, nil
			}
		case '<':
			if l.match('=') {
				return token.Token{Type: token.TokenTypeLessEqual, Lexeme: "<=", Literal: nil, Line: l.line}, nil
			} else {
				return token.Token{Type: token.TokenTypeLess, Lexeme: "<", Literal: nil, Line: l.line}, nil
			}
		case '/':
			if l.match('/') {
				for l.peek() != '\n' && !l.IsAtEnd() {
					l.Advance()
				}

			} else {
				return token.Token{Type: token.TokenTypeSlash, Lexeme: "/", Literal: nil, Line: l.line}, nil
			}
		case ' ':
			noop()
		case '\r':
			noop()
		case '\t':
			noop()
		case '\n':
			l.line++
		case '"':
			str, err := l.nextString()
			if err != nil {
				return token.Token{Type: token.TokenTypeString, Lexeme: str, Literal: str, Line: l.line}, err
			}
			return token.Token{Type: token.TokenTypeString, Lexeme: str, Literal: str, Line: l.line}, nil

		default:
			if isDigit(c) {
				return l.nextNumber()
			} else if isAlpha(c) {
				return l.nextKeywordOrIdentifier()
			}
			return token.Token{Type: token.TokenTypeEOF, Lexeme: string(c), Literal: nil, Line: l.line}, fmt.Errorf("Unexpected character %x", c)

		}
	}

	return token.Token{Type: token.TokenTypeEOF, Lexeme: "", Literal: nil, Line: l.line}, nil
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (l *Lexer) nextKeywordOrIdentifier() (token.Token, error) {
	for isAlpha(l.peek()) || isDigit(l.peek()) {
		l.Advance()
	}

	str := l.source[l.start:l.current]
	switch str {
	case "and":
		return token.Token{Type: token.TokenTypeAnd, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "class":
		return token.Token{Type: token.TokenTypeClass, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "else":
		return token.Token{Type: token.TokenTypeElse, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "false":
		return token.Token{Type: token.TokenTypeFalse, Lexeme: str, Literal: false, Line: l.line}, nil
	case "for":
		return token.Token{Type: token.TokenTypeFor, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "fun":
		return token.Token{Type: token.TokenTypeFun, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "if":
		return token.Token{Type: token.TokenTypeIf, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "nil":
		return token.Token{Type: token.TokenTypeNil, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "or":
		return token.Token{Type: token.TokenTypeOr, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "print":
		return token.Token{Type: token.TokenTypePrint, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "return":
		return token.Token{Type: token.TokenTypeReturn, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "super":
		return token.Token{Type: token.TokenTypeSuper, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "this":
		return token.Token{Type: token.TokenTypeThis, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "true":
		return token.Token{Type: token.TokenTypeTrue, Lexeme: str, Literal: true, Line: l.line}, nil
	case "var":
		return token.Token{Type: token.TokenTypeVar, Lexeme: str, Literal: nil, Line: l.line}, nil
	case "while":
		return token.Token{Type: token.TokenTypeWhile, Lexeme: str, Literal: nil, Line: l.line}, nil
	default:
		return token.Token{Type: token.TokenTypeIdentifier, Lexeme: str, Literal: nil, Line: l.line}, nil
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (l *Lexer) nextNumber() (token.Token, error) {
	for isDigit(l.peek()) {
		l.Advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.Advance()

		for isDigit(l.peek()) {
			l.Advance()
		}
	}

	str := l.source[l.start:l.current]
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return token.Token{Type: token.TokenTypeNumber, Lexeme: str, Literal: nil, Line: l.line}, err
	}
	return token.Token{Type: token.TokenTypeNumber, Lexeme: str, Literal: num, Line: l.line}, nil
}

func (l *Lexer) nextString() (string, error) {
	for l.peek() != '"' && !l.IsAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.Advance()
	}
	if l.IsAtEnd() {
		return "", errors.New("unterminated string.")
	}

	l.Advance()

	str := l.source[l.start+1 : l.current-1]
	return str, nil
}

func noop() {

}
