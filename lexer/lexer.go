package lexer

import (
	"fmt"
	gotoken "go/token"
	"pl0Compiler/token"
)

func PosString(fileName string, src string, pos int) string {
	fSet := gotoken.NewFileSet()
	fSet.AddFile(fileName, 1, len(src)).SetLinesForContent([]byte(src))
	return fmt.Sprintf("%v", fSet.Position(gotoken.Pos(pos+1)))
}

type Lexer struct {
	src      *SourceStream
	tokens   []token.Token
	comments []token.Token
}

func NewLexer(name, input string) *Lexer {
	p := &Lexer{src: NewSourceStream(name, input)}
	p.run()
	return p
}

func (p *Lexer) Tokens() []token.Token {
	return p.tokens
}

func (p *Lexer) Comments() []token.Token {
	return p.comments
}

func (p *Lexer) emit(typ token.TokenType) {
	lit, pos := p.src.EmitToken()
	if typ == token.IDENT {
		typ = token.LoopUp(lit)
	}
	p.tokens = append(p.tokens, token.Token{
		Type:    typ,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) emitComment() {
	lit, pos := p.src.EmitToken()
	p.comments = append(p.comments, token.Token{
		Type:    token.COMMENT,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) errorf(format string, args ...interface{}) {
	tok := token.Token{
		Type:    token.ERROR,
		Literal: fmt.Sprintf(format, args...),
		Pos:     token.Pos(p.src.pos),
	}
	p.tokens = append(p.tokens, tok)
	panic(tok)
}

func (p *Lexer) run() (tokens []token.Token) {
	defer func() {
		tokens = p.tokens
		if r := recover(); r != nil {
			if _, ok := r.(token.Token); !ok {
				panic(r)
			}
		}
	}()

	for {
		r := p.src.Read()
		if r == rune(token.EOF) {
			p.emit(token.EOF)
			return
		}

		switch {
		case r == ';':
			p.emit(token.SEMICOLON)
		case isSpace(r):
			p.src.IgnoreToken()
		case isAlpha(r):
			p.src.Unread()
			for {
				if r := p.src.Read(); !isAlphaNumberic(r) {
					p.src.Unread()
					p.emit(token.IDENT)
					break
				}
			}
		case '0' <= r && r <= '9': // 123, 4.5
			p.src.Unread()
			digits := "0123456789"
			p.src.AcceptRun(digits)
			p.emit(token.NUMBER)
		case r == '+': // +
			p.emit(token.ADD)
		case r == '-': // -
			p.emit(token.SUB)
		case r == '*': // *
			p.emit(token.MUL)
		case r == '/': // /
			p.emit(token.DIV)
			//peek := p.src.Peek()
			//if peek == '/' {
			//	// line comment
			//	for {
			//		t := p.src.Read()
			//		if t == '\n' {
			//			p.src.Unread()
			//			p.emitComment()
			//			break
			//		}
			//		if t == rune(token.EOF) {
			//			p.emitComment()
			//			return
			//		}
			//	}
			//} else {
			//	p.emit(token.DIV)
			//}
		case r == '=': // =
			p.emit(token.EQL)
		case r == '<': // < <= <>
			switch p.src.Read() {
			case '=':
				p.emit(token.LEQ)
			case '>':
				p.emit(token.NEQ)
			default:
				p.src.Unread()
				p.emit(token.LSS)
			}
		case r == '>': // > >=
			switch p.src.Read() {
			case '=':
				p.emit(token.GEQ)
			default:
				p.src.Unread()
				p.emit(token.GTR)
			}
		case r == ':':
			p.src.Read()
			p.emit(token.ASSIGN)
		case r == '.':
			p.emit(token.PERIOD)
		case r == ',':
			p.emit(token.COMMA)
		case r == '(': // (
			p.emit(token.LPAREN)
			//peek := p.src.Read()
			//if peek == '*' {
			//	// multiline comment
			//	for {
			//		t := p.src.Read()
			//		if t == '*' && p.src.Read() == ')' {
			//			p.src.Read()
			//			p.emitComment()
			//			break
			//		}
			//		if t == rune(token.EOF) {
			//			p.errorf("unterminated quoted string")
			//			return
			//		}
			//	}
			//} else {
			//	p.emit(token.LPAREN)
			//}
		case r == ')':
			p.emit(token.RPAREN)
		case r == '{':
			p.src.IgnoreToken()
			for {
				t := p.src.Read()
				if t == '}' {
					p.src.Unread()
					p.emitComment()
					p.src.IgnoreToken()
					break
				}
				if t == rune(token.EOF) || t == '\n' {
					p.errorf("unterminated quoted string")
					return
				}
			}
		}
	}
}

func Lex(name, input string) (tokens, comments []token.Token) {
	l := NewLexer(name, input)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}
