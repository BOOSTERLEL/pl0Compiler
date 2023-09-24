package ast

import "pl0Compiler/token"

func (p *Program) Pos() token.Pos {
	return token.NoPos
}

func (p *Program) End() token.Pos {
	return token.NoPos
}

func (p *Program) nodeType() {

}

func (i Ident) Pos() token.Pos {
	return token.NoPos
}

func (i Ident) End() token.Pos {
	return token.NoPos
}

func (i Ident) exprType() {

}

func (b BinaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (b BinaryExpr) End() token.Pos {
	return token.NoPos
}

func (b BinaryExpr) exprType() {

}

func (c CallExpr) Pos() token.Pos {
	return token.NoPos
}

func (c CallExpr) End() token.Pos {
	return token.NoPos
}

func (c CallExpr) exprType() {

}

func (u UnaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (u UnaryExpr) End() token.Pos {
	return token.NoPos
}

func (u UnaryExpr) exprType() {

}

func (e ExprStmt) Pos() token.Pos {
	return token.NoPos
}

func (e ExprStmt) End() token.Pos {
	return token.NoPos
}

func (e ExprStmt) stmtType() {

}

func (b BlockStmt) Pos() token.Pos {
	return token.NoPos
}

func (b BlockStmt) End() token.Pos {
	return token.NoPos
}

func (b BlockStmt) stmtType() {

}

func (v VarDecl) Pos() token.Pos {
	return token.NoPos
}

func (v VarDecl) End() token.Pos {
	return token.NoPos
}

func (v VarDecl) stmtType() {

}

func (a AssignStmt) Pos() token.Pos {
	return token.NoPos
}

func (a AssignStmt) End() token.Pos {
	return token.NoPos
}

func (a AssignStmt) stmtType() {

}

func (i IfStmt) Pos() token.Pos {
	return token.NoPos
}

func (i IfStmt) End() token.Pos {
	return token.NoPos
}

func (i IfStmt) stmtType() {

}

func (w WhileStmt) Pos() token.Pos {
	return token.NoPos
}

func (w WhileStmt) End() token.Pos {
	return token.NoPos
}

func (w WhileStmt) stmtType() {

}

func (n Number) Pos() token.Pos {
	return token.NoPos
}

func (n Number) End() token.Pos {
	return token.NoPos
}

func (n Number) exprType() {

}

func (c CallExpr) stmtType() {

}

func (i Ident) nodeType() {

}

func (p ParenExpr) exprType() {

}

func (p ProcDecl) Pos() token.Pos {
	return token.NoPos
}

func (p ProcDecl) End() token.Pos {
	return token.NoPos
}

func (p ProcDecl) nodeType() {

}

func (p ParenExpr) Pos() token.Pos {
	return token.NoPos
}

func (p ParenExpr) End() token.Pos {
	return token.NoPos
}

func (v VarDecl) nodeType() {

}

func (d DefineStmt) Pos() token.Pos {
	return token.NoPos
}

func (d DefineStmt) End() token.Pos {
	return token.NoPos
}

func (d DefineStmt) nodeType() {

}

func (r RepeatStmt) Pos() token.Pos {
	return token.NoPos
}

func (r RepeatStmt) End() token.Pos {
	return token.NoPos
}

func (r RepeatStmt) stmtType() {

}

func (I IOStmt) Pos() token.Pos {
	return token.NoPos
}

func (I IOStmt) End() token.Pos {
	return token.NoPos
}

func (I IOStmt) stmtType() {

}
