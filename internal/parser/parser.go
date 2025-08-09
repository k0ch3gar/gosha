package parser

import (
	"fmt"
	"strconv"
	"strings"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/token"
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
	NIL     = &ast.NilDataType{}
	ANY     = &ast.AnyDataType{}
	INT     = &ast.IntegerDataType{}
	STRING  = &ast.StringDataType{}
	BOOLEAN = &ast.BooleanDataType{}
	RETURN  = &ast.ReturnDataType{}
	ERROR   = &ast.ErrorDataType{}
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
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BASH, p.parseBashExpression)

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

func (p *Parser) parseStringLiteral() ast.Expression {
	expression := &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return expression
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
	lit.ReturnType = NIL

	if p.peekTokenIs(token.DTYPE) || p.peekTokenIs(token.FUNCTION) {
		p.nextToken()
		lit.ReturnType = p.parseDataType()
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

	var dataType *ast.DataType = new(ast.DataType)

	ident := &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		DataType: dataType,
	}

	if p.peekTokenIs(token.DTYPE) {
		p.nextToken()
		*dataType = p.parseDataType()
		dataType = new(ast.DataType)
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
			*dataType = p.parseDataType()
			dataType = new(ast.DataType)
		}

		idents = append(idents, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if idents[len(idents)-1].DataType == nil {
		msg := fmt.Sprintf("expected parameter type for %s", idents[len(idents)-1].Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	return idents
}

func (p *Parser) parseIfStatement() ast.Statement {
	expression := &ast.IfStatement{Token: p.curToken}

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
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		DataType: nil,
	}

	if p.peekTokenIs(token.DTYPE) || p.peekTokenIs(token.FUNCTION) {
		p.nextToken()

		dType := p.parseDataType()
		stmt.Name.DataType = &dType
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func rawTypeToDataType(raw string) ast.DataType {
	switch raw {
	case "string":
		return STRING
	case "int":
		return INT
	case "bool":
		return BOOLEAN
	case "any":
		return ANY
	default:
		return NIL
	}
}

func (p *Parser) parseDataType() ast.DataType {
	if p.curTokenIs(token.DTYPE) {
		return rawTypeToDataType(p.curToken.Literal)
	}

	if !p.curTokenIs(token.FUNCTION) {
		return nil
	}

	dType := &ast.FunctionDataType{}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RPAREN) {
		param := p.parseDataType()
		if param == nil {
			return nil
		}

		dType.Parameters = append(dType.Parameters, p.parseDataType())
		p.nextToken()
		p.nextToken()
	}

	if !p.curTokenIs(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.DTYPE) || p.peekTokenIs(token.FUNCTION) {
		p.nextToken()
		returnType := p.parseDataType()
		if returnType == nil {
			dType.ReturnValue = NIL
		} else {
			dType.ReturnValue = returnType
		}
	}

	return dType
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.HASH:
		return p.parseCommentStatement()
	case token.VAR:
		return p.parseVarStatement()
	//case token.CALL:
	//	return p.parseCallStatement()
	case token.IF:
		return p.parseIfStatement()
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

func (p *Parser) parseCommentStatement() ast.Statement {
	if p.curTokenIs(token.HASH) {
		for !p.curTokenIs(token.NLINE) && !p.curTokenIs(token.EOF) {
			p.nextToken()
		}

	}

	return nil
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

func (p *Parser) parseBashExpression() ast.Expression {
	bashExpr := &ast.BashExpression{
		Token: p.curToken,
		Value: strings.Fields(p.curToken.Literal),
	}

	return bashExpr
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
