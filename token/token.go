package token

import "fmt"

type TokenType int

const (
	TokenTypeLeftParen TokenType = iota
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeSlash
	TokenTypeStar
	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeLess
	TokenTypeLessEqual
	TokenTypeIdentifier
	TokenTypeString
	TokenTypeNumber
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFor
	TokenTypeFun
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile
	TokenTypeQuestionMark
	TokenTypeColon
	TokenTypeEOF
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeLeftParen:
		return "LEFT_PAREN"
	case TokenTypeRightParen:
		return "RIGHT_PAREN"
	case TokenTypeLeftBrace:
		return "LEFT_BRACE"
	case TokenTypeRightBrace:
		return "RIGHT_BRACE"
	case TokenTypeComma:
		return "COMMA"
	case TokenTypeDot:
		return "DOT"
	case TokenTypeMinus:
		return "MINUS"
	case TokenTypePlus:
		return "PLUS"
	case TokenTypeSemicolon:
		return "SEMICOLON"
	case TokenTypeSlash:
		return "SLASH"
	case TokenTypeStar:
		return "STAR"
	case TokenTypeBang:
		return "BANG"
	case TokenTypeBangEqual:
		return "BANG_EQUAL"
	case TokenTypeEqual:
		return "EQUAL"
	case TokenTypeEqualEqual:
		return "EQUAL_EQUAL"
	case TokenTypeGreater:
		return "GREATER"
	case TokenTypeGreaterEqual:
		return "GREATER_EQUAL"
	case TokenTypeLess:
		return "LESS"
	case TokenTypeLessEqual:
		return "LESS_EQUAL"
	case TokenTypeIdentifier:
		return "IDENTIFIER"
	case TokenTypeString:
		return "STRING"
	case TokenTypeNumber:
		return "NUMBER"
	case TokenTypeAnd:
		return "AND"
	case TokenTypeClass:
		return "CLASS"
	case TokenTypeElse:
		return "ELSE"
	case TokenTypeFalse:
		return "FALSE"
	case TokenTypeFor:
		return "FOR"
	case TokenTypeFun:
		return "FUN"
	case TokenTypeIf:
		return "IF"
	case TokenTypeNil:
		return "NIL"
	case TokenTypeOr:
		return "OR"
	case TokenTypePrint:
		return "PRINT"
	case TokenTypeReturn:
		return "RETURN"
	case TokenTypeSuper:
		return "SUPER"
	case TokenTypeThis:
		return "THIS"
	case TokenTypeTrue:
		return "TRUE"
	case TokenTypeVar:
		return "VAR"
	case TokenTypeWhile:
		return "WHILE"
	case TokenTypeQuestionMark:
		return "QUESTION_MARK"
	case TokenTypeColon:
		return "COLON"
	case TokenTypeEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) IsTokenType(targetType TokenType) bool {
	return t.Type == targetType
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}
