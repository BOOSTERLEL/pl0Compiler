package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
	"strconv"
)

func (p *Parser) parseExpr() ast.Expr {
	return p.parseExprBinary(1)
}

func (p *Parser) parseExprBinary(prec int) ast.Expr {
	x := p.parseExprUnary()
	for {
		op := p.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}
		p.MustAcceptToken(op.Type)
		y := p.parseExprBinary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{
			OpPos: op.Pos,
			Op:    op.Type,
			X:     x,
			Y:     y,
		}
	}
}

func (p *Parser) parseExprUnary() ast.Expr {
	if _, ok := p.AcceptToken(token.ADD); ok {
		return p.parseExprPrimary()
	}
	if tok, ok := p.AcceptToken(token.SUB, token.ODD); ok {
		return &ast.UnaryExpr{
			OpPos: tok.Pos,
			Op:    tok.Type,
			X:     p.parseExprPrimary(),
		}
	}
	return p.parseExprPrimary()
}

func (p *Parser) parseExprPrimary() ast.Expr {
	if _, ok := p.AcceptToken(token.LPAREN); ok {
		expr := p.parseExpr()
		p.MustAcceptToken(token.RPAREN)
		return expr
	}

	switch tok := p.PeekToken(); tok.Type {
	case token.NUMBER:
		tokNum := p.MustAcceptToken(token.NUMBER)
		value, _ := strconv.Atoi(tokNum.Literal)
		return &ast.Number{
			ValuePos: tokNum.Pos,
			ValueEnd: tokNum.Pos + token.Pos(len(tokNum.Literal)),
			Value:    value,
		}
	case token.IDENT:
		p.MustAcceptToken(token.IDENT)
		return &ast.Ident{
			NamePos: tok.Pos,
			Name:    tok.Literal,
		}
	default:
		p.errorf(tok.Pos, "unknown tok: type=%v, lit=%q", tok.Type, tok.Literal)
		panic("unreachable")
	}

}

func (p *Parser) parseExprCall() *ast.CallStmt {
	tokIdent := p.MustAcceptToken(token.IDENT)
	p.MustAcceptToken(token.SEMICOLON)

	return &ast.CallStmt{
		ProcedureName: &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.Literal,
		},
	}
}
