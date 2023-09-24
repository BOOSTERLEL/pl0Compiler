package token

import (
	"fmt"
	"strconv"
)

// TokenType 词法记号类型
type TokenType int

// Token 记号值
type Token struct {
	Type    TokenType // 记号的类型
	Pos     Pos       // 记号所在的位置(从1开始)
	Literal string    // 程序中原始的字符串
}

// 记号类型
const (
	EOF TokenType = iota
	ERROR
	COMMENT

	IDENT
	NUMBER

	BEGIN
	CALL
	CONST
	DO
	END
	IF
	ODD
	PROCEDURE
	THEN
	VAR
	WHILE
	ELSE
	REPEAT
	UNTIL
	READ
	WRITE

	ADD // +
	SUB // -
	MUL // *
	DIV // /

	EQL // =
	NEQ // <>
	LSS // <
	LEQ // <=
	GTR // >
	GEQ // >=

	ASSIGN //  :=

	LPAREN // (
	RPAREN // )

	COMMA     // ,
	SEMICOLON // ;
	PERIOD    // .
)

func (op TokenType) Precedence() int {
	switch op {
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 1
	case ADD, SUB:
		return 2
	case MUL, DIV:
		return 3
	}
	return 0
}

var tokens = [...]string{
	EOF:     "EOF",
	ERROR:   "ERROR",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",

	BEGIN:     "begin",
	CALL:      "call",
	CONST:     "const",
	DO:        "do",
	END:       "end",
	IF:        "if",
	ODD:       "odd",
	PROCEDURE: "procedure",
	THEN:      "then",
	VAR:       "var",
	WHILE:     "while",
	ELSE:      "else",
	REPEAT:    "repeat",
	UNTIL:     "until",
	READ:      "read",
	WRITE:     "write",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	EQL: "=",
	NEQ: "<>",
	LSS: "<",
	LEQ: "<=",
	GTR: ">",
	GEQ: ">=",

	ASSIGN: ":=",

	LPAREN: "(",
	RPAREN: ")",

	COMMA:     ",",
	SEMICOLON: ";",
	PERIOD:    ".",
}

func (op TokenType) String() string {
	s := ""
	if 0 <= op && op < TokenType(len(tokens)) {
		s = tokens[op]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(op)) + ")"
	}
	return s
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v : \"%v\")  ", t.Type, t.Literal)
}

func (t Token) IntValue() int {
	x, err := strconv.ParseInt(t.Literal, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(x)
}

var keywords = map[string]TokenType{
	"begin":     BEGIN,
	"call":      CALL,
	"const":     CONST,
	"do":        DO,
	"end":       END,
	"if":        IF,
	"odd":       ODD,
	"procedure": PROCEDURE,
	"then":      THEN,
	"var":       VAR,
	"while":     WHILE,
	"else":      ELSE,
	"repeat":    REPEAT,
	"until":     UNTIL,
	"read":      READ,
	"write":     WRITE,
}

func LoopUp(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}
