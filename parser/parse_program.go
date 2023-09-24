package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseProgram() {
	p.program = &ast.Program{}

	p.program.FileName = p.fileName
	p.program.Source = p.src

	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			return
		case token.ERROR:
			panic(tok)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)
		case token.VAR:
			p.program.Globals = append(p.program.Globals, p.parseStmtVar())
		case token.CONST:
			p.program.Const = append(p.program.Const, p.parseStmtConst())
		case token.PROCEDURE:
			p.program.Funcs = append(p.program.Funcs, p.parseProcedure())
		case token.BEGIN:
			p.program.Stmt = p.parseStmtBlock()
			p.MustAcceptToken(token.PERIOD)
		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}
}
