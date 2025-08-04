package parser

import (
	"fmt"
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var (
	VOID = &ast.DataType{Token: token.Token{Type: token.DTYPE, Literal: "void"}, Name: "void"}
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	expression.Arguments = p.parseCallArguments()
	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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
	lit.ReturnType = VOID

	if p.peekTokenIs(token.DTYPE) {
		p.nextToken()
		lit.ReturnType = &ast.DataType{Token: p.curToken, Name: p.curToken.Literal}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var idents []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken()

	dataType := &ast.DataType{Name: token.ILLEGAL}

	ident := &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		DataType: dataType,
	}

	if p.peekTokenIs(token.DTYPE) {
		p.nextToken()
		dataType.Token = p.curToken
		dataType.Name = p.curToken.Literal
		dataType = &ast.DataType{Name: token.ILLEGAL}
	}

	idents = append(idents, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		ident := &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			DataType: dataType,
		}

		if p.peekTokenIs(token.DTYPE) {
			p.nextToken()
			dataType.Token = p.curToken
			dataType.Name = p.curToken.Literal
			dataType = &ast.DataType{Name: token.ILLEGAL}
		}

		idents = append(idents, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if idents[len(idents)-1].DataType.Name == token.ILLEGAL {
		msg := fmt.Sprintf("expected parameter type for %s", idents[len(idents)-1].Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	return idents
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}

	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseBoolean() ast.Expression {
	expression := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.DTYPE) {
		p.nextToken()

		stmt.Name.DataType = &ast.DataType{Token: p.curToken, Name: p.curToken.Literal}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	//case token.CALL:
	//	return p.parseCallStatement()
	case token.SEMICOLON:
		return nil
	case token.NLINE:
		return nil
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		if p.peekTokenIs(token.INITASSIGN) {
			return p.parseInitAssignStatement()
		} else if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		} else {
			return p.parseExpressionStatement()
		}
	default:
		return p.parseExpressionStatement()
	}
}

//func (p *Parser) parseCallStatement()

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseInitAssignStatement() *ast.InitAssignStatement {
	if !p.peekTokenIs(token.INITASSIGN) {
		p.peekError(token.INITASSIGN)
		return nil
	}

	stmt := &ast.InitAssignStatement{
		Token: p.peekToken,
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	//defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NLINE) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	if !p.peekTokenIs(token.ASSIGN) {
		p.peekError(token.ASSIGN)
		return nil
	}

	stmt := &ast.AssignStatement{Token: p.peekToken}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	//defer untrace(trace("parseIntegerLiteral"))
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
