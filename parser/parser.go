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

func (p *Parser) Parse() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0)
	for p.current != len(p.tokens) && !p.currentTokenIs(token.TokenTypeEOF) {
		stmt, err := p.ParseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)

	}

	return statements, nil
}

func (p *Parser) ParseDeclaration() (ast.Stmt, error) {
	if p.currentTokenIs(token.TokenTypeVar) {
		return p.parseVarDeclaration()
	}

	return p.ParseStatement()
}

func (p *Parser) parseVarDeclaration() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeVar) {
		return nil, fmt.Errorf("expected `var` but got token %s", p.currentToken().Type)
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
	}

	// TODO: do synchronize when the parser goes into panic mode.
	if !p.currentTokenIs(token.TokenTypeIdentifier) {
		return nil, fmt.Errorf("expected identifier but got token %s", p.currentToken().Type)
	}
	name, err := p.advance()
	if err != nil {
		return nil, err
	}
	varDeclaration := &ast.VarStatement{
		Name: &name,
	}

	if p.currentTokenIs(token.TokenTypeEqual) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}

		initializer, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		varDeclaration.Initializer = initializer
	}

	_, err = p.consume(token.TokenTypeSemicolon, "expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return varDeclaration, nil
}

func (p *Parser) ParseStatement() (ast.Stmt, error) {
	if p.currentTokenIs(token.TokenTypePrint) {
		return p.parsePrintStatement()
	} else if p.currentTokenIs(token.TokenTypeLeftBrace) {
		return p.parseBlockStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parsePrintStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypePrint) {
		return nil, fmt.Errorf("expected `print` but got token %s", p.currentToken().Type)
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenTypeSemicolon, "expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) parseBlockStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeLeftBrace) {
		return nil, fmt.Errorf("expected `{ }` but got token %s", p.currentToken().Type)
	}

	_, err := p.advance()
	if err != nil {
		return nil, err
	}

	statements := make([]ast.Stmt, 0)
	for !p.currentTokenIs(token.TokenTypeRightBrace) {
		stmt, err := p.ParseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, err = p.advance()
	if err != nil {
		return nil, err
	}

	return &ast.BlockStatement{
		Statements: statements,
	}, nil
}

func (p *Parser) parseExpressionStatement() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenTypeSemicolon, "expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) parseExpression() (ast.Expr, error) {
	return p.parseCommaOperator()
}

func (p *Parser) parseCommaOperator() (ast.Expr, error) {
	left, err := p.parseAssignment()
	if err != nil {
		return nil, err
	}

	for p.currentTokenIs(token.TokenTypeComma) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}
		left = &ast.BeginExpression{
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parseAssignment() (ast.Expr, error) {
	expr, err := p.parseTernary()
	if err != nil {
		return nil, err
	}
	if p.currentTokenIs(token.TokenTypeEqual) {
		_, err = p.advance()
		if err != nil {
			return nil, err
		}

		val, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}

		if variableExpr, ok := expr.(*ast.VariableExpression); ok {
			return &ast.AssignExpression{
				Name:  variableExpr.Name,
				Value: val,
			}, nil
		} else {
			return nil, fmt.Errorf("invalid assignment target %T", expr)

		}

	}

	return expr, nil
}

func (p *Parser) parseTernary() (ast.Expr, error) {
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
	consequent, err := p.parseExpression()
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
	alternative, err := p.parseExpression()
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

func (p *Parser) consume(tokenType token.TokenType, errorMessage string) (token.Token, error) {
	if p.currentTokenIs(tokenType) {
		t, err := p.advance()
		if err != nil {
			return token.Token{}, err
		}
		return t, nil
	} else {
		return token.Token{}, fmt.Errorf(errorMessage)
	}
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
		exp, err := p.parseExpression()
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

	if p.currentTokenIs(token.TokenTypeIdentifier) {
		name, err := p.advance()
		if err != nil {
			return nil, err
		}
		return &ast.VariableExpression{
			Name: name,
		}, nil
	}
	return nil, fmt.Errorf("expected expression got %s", p.currentToken().Type)
}
