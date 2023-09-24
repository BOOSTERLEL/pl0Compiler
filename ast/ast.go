package ast

import (
	"pl0Compiler/token"
)

// Node 表示AST中全部结点
type Node interface {
	Pos() token.Pos
	End() token.Pos
	nodeType()
}

// Program 表示 pl 文件对应的语法树.
type Program struct {
	FileName string // 文件名
	Source   string // 源代码

	Const   []*ConstDecl // 全局常量
	Globals []*VarDecl   // 全局变量
	Funcs   []*ProcDecl  // 函数列表
	Stmt    *BlockStmt   // 程序入口
}

// ConstDecl 常量信息
type ConstDecl struct {
	ConstPos   token.Pos     // const 关键字位置
	Definition []*DefineStmt // 赋值语句
}

// VarDecl 变量信息
type VarDecl struct {
	VarPos token.Pos // var 关键字位置
	Names  []*Ident  // 变量名字
}

// ProcDecl 函数信息
type ProcDecl struct {
	FuncPos token.Pos
	NamePos token.Pos
	Name    string
	VarDecl *VarDecl
	Params  *FieldList
	Body    *BlockStmt
}

// FieldList 参数/属性 列表
type FieldList struct {
	Opening token.Pos
	List    []*Field
	Closing token.Pos
}

// Field 参数/属性
type Field struct {
	Name *Ident
}

// BlockStmt 块语句
type BlockStmt struct {
	BeginPos token.Pos // 'begin'
	List     []Stmt
	EndPos   token.Pos // 'end'
}

type Stmt interface {
	Pos() token.Pos
	End() token.Pos
	stmtType()
}

type ExprStmt struct {
	X Expr
}

// AssignStmt 表示一个赋值语句节点.
type AssignStmt struct {
	Target *Ident    // 要赋值的目标
	OpPos  token.Pos // Op 的位置
	Value  Expr      // 值
}

// DefineStmt 表示一个赋值语句节点.
type DefineStmt struct {
	Target *Ident    // 要赋值的目标
	OpPos  token.Pos // Op 的位置
	Value  Expr      // 值
}

// IfStmt 表示一个 if 语句节点.
type IfStmt struct {
	If   token.Pos  // if 关键字的位置
	Cond Expr       // if 条件, *BinaryExpr
	Body *BlockStmt // if 为真时对应的语句列表
	Else Stmt // else 对应的语句
}

// WhileStmt 表示一个 while 语句节点.
type WhileStmt struct {
	While token.Pos  // while 关键字的位置
	Cond  Expr       // 条件表达式
	Body  *BlockStmt // 循环对应的语句列表
}

// RepeatStmt 表示一个 repeat 语句节点.
type RepeatStmt struct {
	Repeat token.Pos  // repeat 关键字的位置
	Cond   Expr       // 条件表达式
	Body   *BlockStmt // 循环对应的语句列表
}

// IOStmt 表示一个 read/write 语句节点.
type IOStmt struct {
	IOPos  token.Pos       // IO 关键字发位置
	Type   token.TokenType // IO类型
	Params *FieldList      // IO对象
}

type Expr interface {
	Pos() token.Pos
	End() token.Pos
	exprType()
}

// Ident 标识符
type Ident struct {
	NamePos token.Pos
	Name    string
}

// Number 整型
type Number struct {
	ValuePos token.Pos
	ValueEnd token.Pos
	Value    int
}

// BinaryExpr 二元表达式
type BinaryExpr struct {
	OpPos token.Pos       // 运算符位置
	Op    token.TokenType // 运算符类型
	X     Expr            // 左边的运算对象
	Y     Expr            // 右边的运算对象
}

// UnaryExpr 一元表达式
type UnaryExpr struct {
	OpPos token.Pos       // 运算符位置
	Op    token.TokenType // 运算符类型
	X     Expr            // 运算对象
}

// ParenExpr 表示一个圆括弧表达式.
type ParenExpr struct {
	Lparen token.Pos // "(" 的位置
	X      Expr      // 圆括弧内的表达式对象
	Rparen token.Pos // ")" 的位置
}

// CallExpr 表示一个函数调用
type CallExpr struct {
	ProcedureName *Ident    // 函数名字
	CallPos       token.Pos // 位置
	Params        *FieldList
}
