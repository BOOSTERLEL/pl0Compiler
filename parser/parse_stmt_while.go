package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseStmtWhile() *ast.WhileStmt {
	tokFor := p.MustAcceptToken(token.WHILE)

	whileStmt := &ast.WhileStmt{
		While: tokFor.Pos,
	}

	whileStmt.Cond = p.parseExpr()
	p.MustAcceptToken(token.DO)

	if tok := p.PeekToken(); tok.Type == token.BEGIN {
		whileStmt.Body = p.parseStmtBlock()
	} else {
		pos := p.PeekToken().Pos
		stmt := p.parseStmt()
		stmts := make([]ast.Stmt, 1)
		stmts[0] = stmt
		blockStmt := &ast.BlockStmt{
			BeginPos: pos,
			EndPos:   pos,
			List:     stmts,
		}
		whileStmt.Body = blockStmt
	}

	return whileStmt
}

func (p *Parser) parseStmtRepeat() *ast.RepeatStmt {
	tokFor := p.MustAcceptToken(token.REPEAT)

	repeatStmt := &ast.RepeatStmt{
		Repeat: tokFor.Pos,
	}

	if tok := p.PeekToken(); tok.Type == token.BEGIN {
		repeatStmt.Body = p.parseStmtBlock()
	} else {
		pos := p.PeekToken().Pos
		stmt := p.parseStmt()
		stmts := make([]ast.Stmt, 1)
		stmts[0] = stmt
		blockStmt := &ast.BlockStmt{
			BeginPos: pos,
			EndPos:   pos,
			List:     stmts,
		}
		repeatStmt.Body = blockStmt
	}

	p.MustAcceptToken(token.UNTIL)
	repeatStmt.Cond = p.parseExpr()

	return repeatStmt
}
