package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseStmtVar() *ast.VarDecl {
	tokVar := p.MustAcceptToken(token.VAR)
	var varDecl = &ast.VarDecl{
		VarPos: tokVar.Pos,
	}
	for {
		tokArg := p.MustAcceptToken(token.IDENT)
		varDecl.Names = append(varDecl.Names, &ast.Ident{
			NamePos: tokArg.Pos,
			Name:    tokArg.Literal,
		})
		if _, ok := p.AcceptToken(token.SEMICOLON); ok {
			break
		}
		p.MustAcceptToken(token.COMMA)
	}
	return varDecl
}

func (p *Parser) parseStmtConst() *ast.ConstDecl {
	tokConst := p.MustAcceptToken(token.CONST)
	var constDecl = &ast.ConstDecl{
		ConstPos: tokConst.Pos,
	}
	for {
		target := p.parseExprPrimary()
		tok := p.MustAcceptToken(token.EQL)
		expr := p.parseExpr()
		constDecl.Definition = append(constDecl.Definition, &ast.DefineStmt{
			Target: target.(*ast.Ident),
			OpPos:  tok.Pos,
			Value:  expr,
		})
		if _, ok := p.AcceptToken(token.SEMICOLON); ok {
			break
		}
		p.MustAcceptToken(token.COMMA)
	}
	return constDecl
}
