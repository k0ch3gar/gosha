package parser

import (
	"fmt"
	"kstmc.com/gosha/internal/token"
)

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixExpressionParseFuncError(tt token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tt)
	p.errors = append(p.errors, msg)
}
