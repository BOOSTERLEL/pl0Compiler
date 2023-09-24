package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseStmt() ast.Stmt {
	switch tok := p.PeekToken(); tok.Type {
	case token.EOF:
		return nil
	case token.ERROR:
		p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
	case token.SEMICOLON:
		p.AcceptTokenList(token.SEMICOLON)
		return nil
	case token.BEGIN: // begin
		return p.parseStmtBlock()
	case token.VAR:
		return p.parseStmtVar()
	case token.READ, token.WRITE:
		return p.parseIOStmt()
	default:
		return p.parseStmtAssign()
	}
	panic("unreachable")
}

func (p *Parser) parseStmtAssign() ast.Stmt {
	// expr := expr;
	target := p.parseExpr()
	tok := p.MustAcceptToken(token.ASSIGN)
	expr := p.parseExpr()
	p.MustAcceptToken(token.SEMICOLON)
	return &ast.AssignStmt{
		Target: target.(*ast.Ident),
		OpPos:  tok.Pos,
		Value:  expr,
	}
}

func (p *Parser) parseStmtBlock() *ast.BlockStmt {
	block := &ast.BlockStmt{}

	tokBegin := p.MustAcceptToken(token.BEGIN) // begin

Loop:
	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ERROR:
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)
		case token.BEGIN: // begin
			block.List = append(block.List, p.parseStmtBlock())
		case token.END: // end
			break Loop
		case token.VAR:
			block.List = append(block.List, p.parseStmtVar())
		case token.IF:
			block.List = append(block.List, p.parseStmtIf())
		case token.WHILE:
			block.List = append(block.List, p.parseStmtWhile())
		case token.CALL:
			block.List = append(block.List, p.parseCall())
		case token.REPEAT:
			block.List = append(block.List, p.parseStmtRepeat())
		case token.READ, token.WRITE:
			block.List = append(block.List, p.parseIOStmt())
		default:
			block.List = append(block.List, p.parseStmtAssign())
		}
	}

	tokEnd := p.MustAcceptToken(token.END) // end

	block.BeginPos = tokBegin.Pos
	block.EndPos = tokEnd.Pos

	return block
}

func (p *Parser) parseStmtExpr() *ast.ExprStmt {
	return &ast.ExprStmt{
		X: p.parseExpr(),
	}
}
