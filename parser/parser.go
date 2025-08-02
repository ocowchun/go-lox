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
	} else if p.currentTokenIs(token.TokenTypeFun) {
		if !p.currentTokenIs(token.TokenTypeFun) {
			return nil, fmt.Errorf("expected `fun` but got token %s", p.currentToken().Type)
		} else {
			_, err := p.advance()
			if err != nil {
				return nil, err
			}
		}
		return p.parseFunction("function")
	}

	return p.ParseStatement()
}

func (p *Parser) parseFunction(kind string) (ast.Stmt, error) {
	name, err := p.consume(token.TokenTypeIdentifier, fmt.Sprintf("expected %s name", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenTypeLeftParen, fmt.Sprintf("expected `(` after %s name", kind))
	if err != nil {
		return nil, err
	}

	parameters := make([]token.Token, 0)
	for !p.currentTokenIs(token.TokenTypeRightParen) {
		parameter, err := p.consume(token.TokenTypeIdentifier, fmt.Sprintf("expected parameter name for %s", kind))
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, parameter)

		for !p.currentTokenIs(token.TokenTypeRightParen) {
			_, err = p.consume(token.TokenTypeComma, fmt.Sprintf("expected `,` after argument for %s", kind))
			if err != nil {
				return nil, err
			}

			parameter, err := p.consume(token.TokenTypeIdentifier, fmt.Sprintf("expected parameter name for %s", kind))
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameter)
		}
	}

	_, err = p.consume(token.TokenTypeRightParen, fmt.Sprintf("expected `)` after %s parameters", kind))

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStatement{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}, nil
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
	switch p.currentToken().Type {
	case token.TokenTypeIf:
		return p.parseIfStatement()
	case token.TokenTypePrint:
		return p.parsePrintStatement()
	case token.TokenTypeLeftBrace:
		return p.parseBlockStatement()
	case token.TokenTypeWhile:
		return p.parseWhileStatement()
	case token.TokenTypeFor:
		return p.parseForStatement()
	case token.TokenTypeReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeReturn) {
		return nil, fmt.Errorf("expected `return` but got token %s", p.currentToken().Type)
	}
	keyword, err := p.advance()
	if err != nil {
		return nil, err
	}

	var exp ast.Expr
	if !p.currentTokenIs(token.TokenTypeSemicolon) {
		exp, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.TokenTypeSemicolon, "expect `;` after return statement")
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStatement{
		Keyword: keyword,
		Value:   exp,
	}, nil
}

func (p *Parser) parseForStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeFor) {
		return nil, fmt.Errorf("expected `for` but got token %s", p.currentToken().Type)
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(token.TokenTypeLeftParen, "expect '(' after `for`")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if p.currentTokenIs(token.TokenTypeSemicolon) {
		_, err = p.advance()
		if err != nil {
			return nil, err
		}
	} else if p.currentTokenIs(token.TokenTypeVar) {
		initializer, err = p.parseVarDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.parseExpressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr
	if !p.currentTokenIs(token.TokenTypeSemicolon) {
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenTypeSemicolon, "expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.currentTokenIs(token.TokenTypeRightParen) {
		increment, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenTypeRightParen, "expect ')' after loop clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.ParseStatement()

	if increment != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Stmt{
				body,
				&ast.ExpressionStatement{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &ast.LiteralExpression{Value: true}
	}
	body = &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) parseWhileStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeWhile) {
		return nil, fmt.Errorf("expected `while` but got token %s", p.currentToken().Type)
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(token.TokenTypeLeftParen, "expect '(' after `while`")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenTypeRightParen, "expect ')' after `while` condition")
	if err != nil {
		return nil, err
	}

	body, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parseIfStatement() (ast.Stmt, error) {
	if !p.currentTokenIs(token.TokenTypeIf) {
		return nil, fmt.Errorf("expected `if` but got token %s", p.currentToken().Type)
	} else {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(token.TokenTypeLeftParen, "expect '(' after `if`")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenTypeRightParen, "expect ')' after `if` condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.ParseStatement()
	if err != nil {
		return nil, err
	}

	if p.currentTokenIs(token.TokenTypeElse) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}

		elseBranch, err := p.ParseStatement()
		if err != nil {
			return nil, err
		}

		return &ast.IfStatement{
			Condition:  condition,
			ThenBranch: thenBranch,
			ElseBranch: elseBranch,
		}, nil
	}

	return &ast.IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
	}, nil
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

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	if !p.currentTokenIs(token.TokenTypeLeftBrace) {
		return nil, fmt.Errorf("expected `{` but got token %s", p.currentToken().Type)
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
	return p.parseCommaExpression()
}

func (p *Parser) parseCommaExpression() (ast.Expr, error) {
	expr, err := p.parseAssignment()
	if err != nil {
		return nil, err
	}

	if !p.currentTokenIs(token.TokenTypeComma) {
		return expr, nil
	}

	expressions := []ast.Expr{expr}

	for p.currentTokenIs(token.TokenTypeComma) {
		_, err = p.advance()
		if err != nil {
			return nil, err
		}

		expr, err = p.parseAssignment()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return &ast.CommaExpression{
		Expressions: expressions,
	}, nil
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
	predicate, err := p.ParseOr()
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

func (p *Parser) ParseOr() (ast.Expr, error) {
	expr, err := p.ParseAnd()
	if err != nil {
		return nil, err
	}
	for p.currentTokenIs(token.TokenTypeOr) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.ParseAnd()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) ParseAnd() (ast.Expr, error) {
	expr, err := p.ParseEquality()
	if err != nil {
		return nil, err
	}
	for p.currentTokenIs(token.TokenTypeAnd) {
		op, err := p.advance()
		if err != nil {
			return nil, err
		}

		right, err := p.ParseEquality()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
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
		return token.Token{}, fmt.Errorf("%s got token %s", errorMessage, p.currentToken().Lexeme)
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

	return p.parseCall()
}

func (p *Parser) parseCall() (ast.Expr, error) {
	callee, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if p.currentTokenIs(token.TokenTypeLeftParen) {
			_, err := p.advance()
			if err != nil {
				return nil, err
			}
			callee, err = p.finishCall(callee)
		} else {
			break
		}
	}

	return callee, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := make([]ast.Expr, 0)

	if !p.currentTokenIs(token.TokenTypeRightParen) {
		commaExpression, err := p.parseCommaExpression()
		if err != nil {
			return nil, err
		}

		if commaExpression, ok := commaExpression.(*ast.CommaExpression); ok {
			if len(commaExpression.Expressions) >= 255 {
				// TODO: might still want to parse the expression since the syntax is valid.
				return nil, fmt.Errorf("can't have more than 255 arguments., got %d", len(commaExpression.Expressions))
			}

			arguments = append(arguments, commaExpression.Expressions...)
		} else {
			arguments = append(arguments, commaExpression)
		}
	}

	paren, err := p.consume(token.TokenTypeRightParen, "expect `)` after function arguments")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpression{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil

}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	if p.currentTokenIs(token.TokenTypeTrue) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
		return &ast.LiteralExpression{Value: true}, nil
	}

	if p.currentTokenIs(token.TokenTypeFalse) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
		return &ast.LiteralExpression{Value: false}, nil
	}

	if p.currentTokenIs(token.TokenTypeNil) {
		_, err := p.advance()
		if err != nil {
			return nil, err
		}
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
		_, err := p.advance()
		if err != nil {
			return nil, err
		}

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
