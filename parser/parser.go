package parser

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ocowchun/go-lox/ast"
	"github.com/ocowchun/go-lox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.ParseExpression()
}

func (p *Parser) ParseExpression() (ast.Expr, error) {
	return p.ParseCommaOperator()
}

func (p *Parser) ParseCommaOperator() (ast.Expr, error) {
	left, err := p.ParseTernary()
	if err != nil {
		return nil, err
	}
	if !p.currentTokenIs(token.TokenTypeComma) {
		return left, nil
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.ParseCommaOperator()
		if err != nil {
			return nil, err
		}
		return &ast.BeginExpression{
			Left:  left,
			Right: right,
		}, nil
	}
}

func (p *Parser) ParseTernary() (ast.Expr, error) {
	// predicate ? exp1 : exp2
	predicate, err := p.ParseEquality()
	if err != nil {
		return nil, err
	}

	if !p.currentTokenIs(token.TokenTypeQuestionMark) {
		return predicate, nil
	}

	_, err = p.advance()
	if err != nil {
		return nil, err
	}
	consequent, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if !p.currentTokenIs(token.TokenTypeColon) {
		return nil, fmt.Errorf("expected `:` but got token %s", p.currentToken().Type)
	}

	_, err = p.advance()
	if err != nil {
		return nil, err
	}
	alternative, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.ConditionExpression{
		Predicate:   predicate,
		Consequent:  consequent,
		Alternative: alternative,
	}, nil

}

func (p *Parser) ParseEquality() (ast.Expr, error) {
	var left ast.Expr
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.currentTokenIs(token.TokenTypeBangEqual) || p.currentTokenIs(token.TokenTypeEqualEqual) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}

	}

	return left, nil
}

func (p *Parser) currentToken() token.Token {
	if p.current >= len(p.tokens) {
		return token.Token{
			Type: token.TokenTypeEOF,
		}
	}

	return p.tokens[p.current]
}

func (p *Parser) currentTokenIs(tokenTypes ...token.TokenType) bool {
	if p.current >= len(p.tokens) {
		return false
	}

	return slices.Contains(tokenTypes, p.currentToken().Type)
}

func (p *Parser) advance() (token.Token, error) {
	if p.current >= len(p.tokens) {
		return token.Token{}, errors.New("unexpected end of input")
	}

	t := p.tokens[p.current]
	p.current++
	return t, nil
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	var left ast.Expr
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.currentTokenIs(token.TokenTypeGreater, token.TokenTypeGreaterEqual, token.TokenTypeLess, token.TokenTypeLessEqual) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}

	}

	return left, nil
}

func (p *Parser) parseTerm() (ast.Expr, error) {
	var left ast.Expr
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.currentTokenIs(token.TokenTypePlus, token.TokenTypeMinus) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseFactor() (ast.Expr, error) {
	var left ast.Expr
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.currentTokenIs(token.TokenTypeStar, token.TokenTypeSlash) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	if p.currentTokenIs(token.TokenTypeMinus, token.TokenTypeBang) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpression{
			Operator: op,
			Right:    right,
		}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	if p.currentTokenIs(token.TokenTypeTrue) {
		p.advance()
		return &ast.LiteralExpression{Value: true}, nil
	}

	if p.currentTokenIs(token.TokenTypeFalse) {
		p.advance()
		return &ast.LiteralExpression{Value: false}, nil
	}

	if p.currentTokenIs(token.TokenTypeNil) {
		p.advance()
		return &ast.LiteralExpression{Value: nil}, nil
	}

	if p.currentTokenIs(token.TokenTypeNumber, token.TokenTypeString) {
		t, err := p.advance()
		if err != nil {
			return nil, err

		}

		return &ast.LiteralExpression{Value: t.Literal}, nil
	}

	if p.currentTokenIs(token.TokenTypeLeftParen) {
		p.advance()
		exp, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}

		if p.currentTokenIs(token.TokenTypeRightParen) {
			_, err := p.advance()
			if err != nil {
				return nil, err
			}

			return &ast.GroupingExpression{Expression: exp}, nil
		} else {

			return nil, fmt.Errorf("expected `)` but got token %s", p.currentToken().Type)
		}
	}

	return nil, fmt.Errorf("expected expression got %s", p.currentToken().Type)
}
