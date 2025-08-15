package parser

import (
	"fmt"

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
	BUILTIN = &ast.BuiltinDataType{}
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
	token.PERCENT:  SUM,
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
	p.registerPrefix(token.FOPER, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
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
	p.registerInfix(token.PERCENT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
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

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var idents []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken()

	var dataType = new(ast.DataType)

	ident := &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		DataType: dataType,
	}

	if !p.peekTokenIs(token.COMMA) {
		p.nextToken()
		*dataType = p.parseDataTypeLiteral()
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
			*dataType = p.parseDataTypeLiteral()
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

	if !p.peekTokenIs(token.ASSIGN) {
		p.nextToken()

		dType := p.parseDataTypeLiteral()
		stmt.Name.DataType = &dType
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	p.nextToken()

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.HASH:
		return p.parseCommentStatement()
	case token.VAR:
		return p.parseVarStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.IF:
		return p.parseIfStatement()
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

func (p *Parser) parseForStatement() ast.Statement {
	forStmt := &ast.ForStatement{
		Token: p.curToken,
	}

	p.nextToken()
	forStmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	forStmt.Consequence = p.parseBlockStatement()
	p.nextToken()
	return forStmt
}

func (p *Parser) parseCommentStatement() ast.Statement {
	if p.curTokenIs(token.HASH) {
		for !p.curTokenIs(token.NLINE) && !p.curTokenIs(token.EOF) {
			p.nextToken()
		}

	}

	return nil
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

	if p.peekTokenIs(token.NLINE) {
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

func (p *Parser) parseFunctionDataType() ast.DataType {
	dType := &ast.FunctionDataType{}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		param := p.parseDataTypeLiteral()
		if param == nil {
			return nil
		}

		dType.Parameters = append(dType.Parameters, param)
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := p.parseDataTypeLiteral()
		if param == nil {
			return nil
		}

		dType.Parameters = append(dType.Parameters, param)
		p.nextToken()
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		returnType := p.parseDataTypeLiteral()
		if returnType == nil {
			dType.ReturnType = NIL
		} else {
			dType.ReturnType = returnType
		}
	}

	return dType
}

func (p *Parser) parseSliceDataType() ast.DataType {
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	p.nextToken()
	arrayDataType := &ast.SliceDataType{}
	arrayDataType.Type = p.parseDataTypeLiteral()
	if arrayDataType.Type == nil {
		return nil
	}

	return arrayDataType
}
