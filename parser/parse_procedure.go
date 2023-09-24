package parser

import (
	"pl0Compiler/ast"
	"pl0Compiler/token"
)

func (p *Parser) parseProcedure() *ast.ProcDecl {
	tokFunc := p.MustAcceptToken(token.PROCEDURE)
	tokFuncIdent := p.MustAcceptToken(token.IDENT)

	proc := &ast.ProcDecl{
		FuncPos: tokFunc.Pos,
		NamePos: tokFuncIdent.Pos,
		Name:    tokFuncIdent.Literal,
		Params:  &ast.FieldList{},
	}
	if _, ok := p.AcceptToken(token.LPAREN); ok {
		for {
			// args
			tokArg := p.MustAcceptToken(token.IDENT)
			proc.Params.List = append(proc.Params.List, &ast.Field{
				Name: &ast.Ident{
					NamePos: tokArg.Pos,
					Name:    tokArg.Literal,
				},
			})
			// )
			if _, ok := p.AcceptToken(token.RPAREN); ok {
				break
			}
			p.MustAcceptToken(token.COMMA)
		}
	}
	p.MustAcceptToken(token.SEMICOLON)

	if _, ok := p.AcceptToken(token.VAR); ok {
		p.UnreadToken()
		proc.VarDecl = p.parseStmtVar()
	}

	// body: begin ... end
	if _, ok := p.AcceptToken(token.BEGIN); ok {
		p.UnreadToken()
		proc.Body = p.parseStmtBlock()
	}
	p.MustAcceptToken(token.SEMICOLON)

	return proc
}

func (p *Parser) parseCall() *ast.CallStmt {
	tokCall := p.MustAcceptToken(token.CALL)
	tokProcIdent := p.MustAcceptToken(token.IDENT)
	call := &ast.CallStmt{
		ProcedureName: &ast.Ident{
			NamePos: tokProcIdent.Pos,
			Name:    tokProcIdent.Literal,
		},
		CallPos: tokCall.Pos,
	}

	if tokLparen, ok := p.AcceptToken(token.LPAREN); ok {
		call.Lparen = tokLparen.Pos
		for {
			call.Args = append(call.Args, p.parseExpr())
			if tokRparen, ok := p.AcceptToken(token.RPAREN); ok {
				call.Rparen = tokRparen.Pos
				break
			}
			p.MustAcceptToken(token.COMMA)
		}
	}

	p.MustAcceptToken(token.SEMICOLON)
	return call
}

func (p *Parser) parseIOStmt() *ast.IOStmt {
	IOTok := token.Token{}
	params := ast.FieldList{}
	switch tok := p.PeekToken(); tok.Type {
	case token.WRITE:
		IOTok = p.MustAcceptToken(token.WRITE)
		paramTok := p.MustAcceptToken(token.IDENT)
		params.List = append(params.List, &ast.Field{
			Name: &ast.Ident{
				NamePos: paramTok.Pos,
				Name:    paramTok.Literal,
			},
		})
	case token.READ:
		IOTok = p.MustAcceptToken(token.READ)
		for {
			paramTok := p.MustAcceptToken(token.IDENT)
			params.List = append(params.List, &ast.Field{
				Name: &ast.Ident{
					NamePos: paramTok.Pos,
					Name:    paramTok.Literal,
				},
			})
			if _, ok := p.AcceptToken(token.SEMICOLON); ok {
				break
			}
			p.MustAcceptToken(token.COMMA)
		}
	default:
		p.errorf(tok.Pos, "unknown token: %v", tok)
	}
	return &ast.IOStmt{
		IOPos:  IOTok.Pos,
		Type:   IOTok.Type,
		Params: &params,
	}
}
