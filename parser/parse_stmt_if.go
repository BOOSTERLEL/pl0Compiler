package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseStmtIf() *ast.IfStmt {
	tokIf := p.MustAcceptToken(token.IF)

	ifStmt := &ast.IfStmt{
		If: tokIf.Pos,
	}

	ifStmt.Cond = p.parseExpr()
	p.MustAcceptToken(token.THEN)
	if tok := p.PeekToken(); tok.Type == token.BEGIN {
		ifStmt.Body = p.parseStmtBlock()
	} else {
		stmt := p.parseStmt()
		ifStmt.Body = &ast.BlockStmt{}
		ifStmt.Body.List = append(ifStmt.Body.List, stmt)
	}

	if _, ok := p.AcceptToken(token.ELSE); ok {
		switch p.PeekToken().Type {
		case token.IF: // else if
			ifStmt.Else = p.parseStmtIf()
		default:
			if tok := p.PeekToken(); tok.Type == token.BEGIN {
				ifStmt.Else = p.parseStmtBlock()
			} else {
				stmt := p.parseStmt()
				ifStmt.Else = &ast.BlockStmt{}
				ifStmt.Else.(*ast.BlockStmt).List = append(ifStmt.Else.(*ast.BlockStmt).List, stmt)
			}
			//ifStmt.Else = p.parseStmtBlock()
		}
	}

	return ifStmt
}
