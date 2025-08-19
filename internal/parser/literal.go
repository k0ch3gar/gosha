package parser

import (
	"fmt"
	"strconv"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/token"
)

func (p *Parser) parseDataTypeLiteral() ast.DataType {
	switch p.curToken.Type {
	case token.ASTERISK:
		p.nextToken()
		return &ast.ReferenceDataType{ValueType: p.parseDataTypeLiteral()}
	case token.DTYPE:
		return tokenTypeToDataType(p.curToken.Literal)
	case token.FUNCTION:
		return p.parseFunctionDataType()
	case token.LBRACKET:
		return p.parseSliceDataType()
	default:
		msg := fmt.Sprintf("unknown data type that starts with %s token type", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()
	lit.ReturnType = NIL

	if !p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		lit.ReturnType = p.parseDataTypeLiteral()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	expression := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	expression := &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return expression
}
