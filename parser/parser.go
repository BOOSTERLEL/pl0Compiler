package parser

import (
	"fmt"
	"pl0Compiler/ast"
	"pl0Compiler/lexer"
	"pl0Compiler/token"
)

type Parser struct {
	fileName string
	src      string

	*TokenStream
	program *ast.Program
	err     error
}

func (p *Parser) errorf(pos token.Pos, format string, args ...interface{}) {
	p.err = fmt.Errorf("%s: %s", lexer.PosString(p.fileName, p.src, int(pos)), fmt.Sprintf(format, args...))
	panic(p.err)
}

func (p *Parser) ParseProgram() (file *ast.Program, err error) {
	defer func() {
		if r := recover(); r != p.err {
			panic(r)
		}
		file, err = p.program, p.err
	}()

	tokens, comments := lexer.Lex(p.fileName, p.src)
	for _, tok := range tokens {
		if tok.Type == token.ERROR {
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		}
	}

	p.TokenStream = NewTokenStream(p.fileName, p.src, tokens, comments)
	p.parseProgram()
	return
}

func NewParser(fileName, src string) *Parser {
	return &Parser{
		fileName: fileName,
		src:      src,
	}
}

func ParseFile(fileName, src string) (*ast.Program, error) {
	p := NewParser(fileName, src)
	return p.ParseProgram()
}
