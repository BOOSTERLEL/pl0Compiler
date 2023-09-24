package compiler

import (
	"bytes"
	"fmt"
	"io"
	"pl0Compiler/ast"
	"pl0Compiler/builtin"
	"pl0Compiler/token"
)

type Compiler struct {
	program *ast.Program
	scope   *Scope
	nextId  int
}

func NewCompiler() *Compiler {
	return &Compiler{
		scope: NewScope(Universe),
	}
}

func (p *Compiler) Compile(program *ast.Program) string {
	var buf bytes.Buffer

	p.program = program

	p.genHeader(&buf, program)
	p.compileProgram(&buf, program)

	return buf.String()
}

func (p *Compiler) enterScope() {
	p.scope = NewScope(p.scope)
}

func (p *Compiler) leaveScope() {
	p.scope = p.scope.Outer
}

func (p *Compiler) restoreScope(scope *Scope) {
	p.scope = scope
}

func (p *Compiler) genHeader(w io.Writer, program *ast.Program) {
	_, _ = fmt.Fprintf(w, "; program name %s\n", program.FileName)
	_, _ = fmt.Fprint(w, builtin.Header)
}

func (p *Compiler) genMain(w io.Writer, program *ast.Program) {
	_, _ = fmt.Fprintf(w, "define i32 @pl_0_main() {\n")
	p.compileStmt(w, program.Stmt)
	_, _ = fmt.Fprintf(w, "\tret i32 0\n}\n")
	_, _ = fmt.Fprintf(w, builtin.MainMain)
}

func (p *Compiler) compileProgram(w io.Writer, program *ast.Program) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	for _, g := range program.Globals {
		for _, name := range g.Names {
			var mangledName = fmt.Sprintf("@pl_0_%s", name.Name)
			p.scope.Insert(&Object{
				Name:        name.Name,
				MangledName: mangledName,
				Node:        name,
			})
			_, _ = fmt.Fprintf(w, "%s = dso_local global i32 0, align 4\n", mangledName)
		}
	}
	if len(program.Globals) != 0 {
		_, _ = fmt.Fprintln(w)
	}

	for _, c := range program.Const {
		for _, name := range c.Definition {
			var mangledName = fmt.Sprintf("@pl_0_%s", name.Target.Name)
			p.scope.Insert(&Object{
				Name:        name.Target.Name,
				MangledName: mangledName,
				Node:        name,
			})
			_, _ = fmt.Fprintf(w, "%s = dso_local constant i32 %d, align 4\n",
				mangledName, name.Value.(*ast.Number).Value)
		}
	}
	if len(program.Const) != 0 {
		_, _ = fmt.Fprintln(w)
	}

	for _, fn := range program.Funcs {
		var mangledName = fmt.Sprintf("@pl_0_%s", fn.Name)
		p.scope.Insert(&Object{
			Name:        fn.Name,
			MangledName: mangledName,
			Node:        fn,
		})
	}

	for _, fn := range program.Funcs {
		p.compileProcedure(w, fn)
	}

	p.genMain(w, program)
}

func (p *Compiler) compileProcedure(w io.Writer, fn *ast.ProcDecl) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	// args
	var argNameList []string
	for _, arg := range fn.Params.List {
		var mangledName = fmt.Sprintf("%%local_%s.pos.%d", arg.Name.Name, arg.Name.NamePos)
		argNameList = append(argNameList, mangledName)
	}

	if fn.Body == nil {
		_, _ = fmt.Fprintf(w, "declare i32 @pl_0_%s()\n", fn.Name)
		return
	}

	_, _ = fmt.Fprintf(w, "define i32 @pl_0_%s(", fn.Name)
	var first = true
	for i, argRegName := range argNameList {
		if first {
			first = false
			_, _ = fmt.Fprintf(w, "i32 noundef %s.arg%d", argRegName, i)
			continue
		}
		_, _ = fmt.Fprintf(w, ", i32 noundef %s.arg%d", argRegName, i)
	}
	_, _ = fmt.Fprintf(w, ") {\n")

	// proc body
	func() {
		// args+body scope
		defer p.restoreScope(p.scope)
		p.enterScope()

		// args
		for i, arg := range fn.Params.List {
			var argRegName = fmt.Sprintf("%s.arg%d", argNameList[i], i)
			var mangledName = argNameList[i]
			p.scope.Insert(&Object{
				Name:        arg.Name.Name,
				MangledName: mangledName,
				Node:        fn,
			})

			_, _ = fmt.Fprintf(w, "\t%s = alloca i32, align 4\n", mangledName)
			_, _ = fmt.Fprintf(w, "\tstore i32 %s, i32* %s\n", argRegName, mangledName)
		}

		// body
		for _, x := range fn.Body.List {
			p.compileStmt(w, x)
		}
	}()

	//p.compileStmt(w, fn.Body)
	_, _ = fmt.Fprintln(w, "\tret i32 0")
	_, _ = fmt.Fprintln(w, "}")
}

func (p *Compiler) compileStmt(w io.Writer, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.VarDecl:
		for _, name := range stmt.Names {

			var mangledName = fmt.Sprintf("%%local_%s.pos.%d", name.Name, stmt.VarPos)
			p.scope.Insert(&Object{
				Name:        name.Name,
				MangledName: mangledName,
				Node:        stmt,
			})

			_, _ = fmt.Fprintf(w, "\t%s = alloca i32, align 4\n", mangledName)
			_, _ = fmt.Fprintf(w, "\tstore i32 0, i32* %s\n", mangledName)
		}

	case *ast.AssignStmt:
		p.compileStmtAssign(w, stmt)
	case *ast.IfStmt:
		p.compileStmtIf(w, stmt)
	case *ast.WhileStmt:
		p.compileStmtWhile(w, stmt)
	case *ast.RepeatStmt:
		p.compileStmtRepeat(w, stmt)
	case *ast.BlockStmt:
		defer p.restoreScope(p.scope)
		p.enterScope()

		for _, x := range stmt.List {
			p.compileStmt(w, x)
		}
	case *ast.ExprStmt:
		p.compileExpr(w, stmt.X)
	case *ast.CallStmt:
		p.compileStmtCall(w, stmt)
	case *ast.IOStmt:
		p.compileIOStmt(w, stmt)

	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", stmt))
	}
}

func (p *Compiler) compileStmtAssign(w io.Writer, stmt *ast.AssignStmt) {
	var name string
	valueName := p.compileExpr(w, stmt.Value)
	if _, obj := p.scope.Lookup(stmt.Target.Name); obj != nil {
		name = obj.MangledName
	} else {
		panic(fmt.Sprintf("var %s undefined", stmt.Target.Name))
	}
	_, _ = fmt.Fprintf(
		w, "\tstore i32 %s, i32* %s\n",
		valueName, name,
	)
}

func (p *Compiler) compileStmtIf(w io.Writer, stmt *ast.IfStmt) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	ifPos := fmt.Sprintf("%d", p.posLine(stmt.If))
	ifCond := p.genLabelId("if.cond.line" + ifPos)
	ifBody := p.genLabelId("if.body.line" + ifPos)
	ifElse := p.genLabelId("if.else.line" + ifPos)
	ifEnd := p.genLabelId("if.end.line" + ifPos)

	func() {
		defer p.restoreScope(p.scope)
		p.enterScope()

		_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifCond)

		// if.cond
		{
			_, _ = fmt.Fprintf(w, "\n%s:\n", ifCond)
			condValue := p.compileExpr(w, stmt.Cond)
			if stmt.Else != nil {
				_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifElse)
			} else {
				_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifEnd)
			}
		}

		// if.body
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", ifBody)
			p.compileStmt(w, stmt.Body)
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifEnd)
		}()

		// if.else
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", ifElse)
			if stmt.Else != nil {
				p.compileStmt(w, stmt.Else)
				_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifEnd)
			} else {
				_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifEnd)
			}
		}()
	}()

	// end
	_, _ = fmt.Fprintf(w, "\n%s:\n", ifEnd)
}

func (p *Compiler) compileStmtWhile(w io.Writer, stmt *ast.WhileStmt) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	whilePos := fmt.Sprintf("%d", p.posLine(stmt.While))
	whileCond := p.genLabelId("while.cond.line" + whilePos)
	whileBody := p.genLabelId("while.body.line" + whilePos)
	whileEnd := p.genLabelId("while.end.line" + whilePos)

	func() {
		defer p.restoreScope(p.scope)
		p.enterScope()

		_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", whileCond)

		// while.cond
		_, _ = fmt.Fprintf(w, "\n%s:\n", whileCond)
		condValue := p.compileExpr(w, stmt.Cond)
		_, _ = fmt.Fprintf(w, "\tbr i1 %s , label %%%s, label %%%s\n", condValue, whileBody, whileEnd)

		// while.body
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", whileBody)
			p.compileStmt(w, stmt.Body)
		}()
	}()

	// end
	_, _ = fmt.Fprintf(w, "\n%s:\n", whileEnd)
}

func (p *Compiler) compileStmtRepeat(w io.Writer, stmt *ast.RepeatStmt) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	repeatPos := fmt.Sprintf("%d", p.posLine(stmt.Repeat))
	repeatCond := p.genLabelId("repeat.cond.line" + repeatPos)
	repeatBody := p.genLabelId("repeat.body.line" + repeatPos)
	repeatEnd := p.genLabelId("repeat.end.line" + repeatPos)

	func() {
		defer p.restoreScope(p.scope)
		p.enterScope()

		_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", repeatBody)

		// repeat.body
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", repeatBody)
			p.compileStmt(w, stmt.Body)
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", repeatCond)
		}()

		// repeat.cond
		_, _ = fmt.Fprintf(w, "\n%s:\n", repeatCond)
		condValue := p.compileExpr(w, stmt.Cond)
		_, _ = fmt.Fprintf(w, "\tbr i1 %s , label %%%s, label %%%s\n", condValue, repeatEnd, repeatBody)
	}()

	// end
	_, _ = fmt.Fprintf(w, "\n%s:\n", repeatEnd)
}

func (p *Compiler) compileStmtCall(w io.Writer, expr *ast.CallStmt) {
	var fnName string
	if _, obj := p.scope.Lookup(expr.ProcedureName.Name); obj != nil {
		fnName = obj.MangledName
	} else {
		panic(fmt.Sprintf("proc %s undefined", expr.ProcedureName.Name))
	}

	var localNames []string
	for _, arg := range expr.Args {
		localNames = append(localNames, p.compileExpr(w, arg))
	}
	_, _ = fmt.Fprintf(w, "\tcall i32 %s(", fnName)
	first := true
	for _, localName := range localNames {
		if first {
			_, _ = fmt.Fprintf(w, "i32 noundef %s", localName)
			first = false
			continue
		}
		_, _ = fmt.Fprintf(w, ", i32 noundef %s", localName)
	}
	_, _ = fmt.Fprintf(w, ")\n")
}

func (p *Compiler) compileIOStmt(w io.Writer, stmt *ast.IOStmt) {
	var targetName string
	switch stmt.Type {
	case token.READ:
		for _, param := range stmt.Params.List {
			localName := p.compileExpr(w, param.Name)
			_, _ = fmt.Fprintf(w, "\tcall i32 @pl_0_builtin_println(i32 %s)\n",
				localName)
		}
	case token.WRITE:
		if _, obj := p.scope.Lookup(stmt.Params.List[0].Name.Name); obj != nil {
			targetName = obj.MangledName
		} else {
			panic(fmt.Sprintf("var %s undefined", stmt.Params.List[0].Name.Name))
		}
		localName := p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = call i32 @pl_0_builtin_write()\n", localName)
		_, _ = fmt.Fprintf(w, "\tstore i32 %s, i32* %s, align 4\n",
			localName, targetName)
	}
}

func (p *Compiler) compileExpr(w io.Writer, expr ast.Expr) (localName string) {
	switch expr := expr.(type) {
	case *ast.Ident:
		var varName string
		if _, obj := p.scope.Lookup(expr.Name); obj != nil {
			varName = obj.MangledName
		} else {
			panic(fmt.Sprintf("var %s undefined", expr.Name))
		}

		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = load i32, i32* %s, align 4\n",
			localName, varName,
		)
		return localName
	case *ast.Number:
		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
			localName, "add", `0`, expr.Value,
		)
		return localName
	case *ast.BinaryExpr:
		localName = p.genId()
		switch expr.Op {
		case token.ADD:
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "add", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.SUB:
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "sub", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.MUL:
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "mul", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.DIV:
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "div", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName

		case token.EQL: // =
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp eq", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.NEQ: // <>
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp ne", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.LSS: // <
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp slt", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.LEQ: // <=
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sle", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.GTR: // >
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sgt", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.GEQ: // >=
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sge", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		default:
			panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
		}
	case *ast.UnaryExpr:
		if expr.Op == token.SUB {
			localName = p.genId()
			_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "sub", `0`, p.compileExpr(w, expr.X),
			)
			return localName
		}
		return p.compileExpr(w, expr.X)
	case *ast.ParenExpr:
		return p.compileExpr(w, expr.X)

	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
	}
}

func (p *Compiler) posLine(pos token.Pos) int {
	if p.program != nil && p.program.Source != "" {
		line := pos.Position(p.program.FileName, p.program.Source).Line
		return line
	}
	return 0
}

func (p *Compiler) genId() string {
	id := fmt.Sprintf("%%t%d", p.nextId)
	p.nextId++
	return id
}

func (p *Compiler) genLabelId(name string) string {
	id := fmt.Sprintf("%s.%d", name, p.nextId)
	p.nextId++
	return id
}
