package compiler

import (
	"bytes"
	"fmt"
	"io"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/builtin"
	"github.com/chai2010/ugo/token"
)

type Compiler struct {
	file   *ast.File
	scope  *Scope
	nextId int
}

func NewCompiler() *Compiler {
	return &Compiler{
		scope: NewScope(Universe),
	}
}

func (p *Compiler) Compile(file *ast.File) string {
	var buf bytes.Buffer

	p.file = file

	p.genHeader(&buf, file)
	p.compileFile(&buf, file)
	p.genMain(&buf, file)

	return buf.String()
}

func (p *Compiler) enterScope() {
	p.scope = NewScope(p.scope)
}

func (p *Compiler) leaveScope() {
	p.scope = p.scope.Outer
}

func (p *Compiler) genHeader(w io.Writer, file *ast.File) {
	fmt.Fprintf(w, "; package %s\n", file.Pkg.Name)
	fmt.Fprint(w, builtin.Header)
}

func (p *Compiler) genMain(w io.Writer, file *ast.File) {
	if file.Pkg.Name != "main" {
		return
	}
	for _, fn := range file.Funcs {
		if fn.Name == "main" {
			fmt.Fprint(w, builtin.MainMain)
			return
		}
	}
}

func (p *Compiler) genInit(w io.Writer, file *ast.File) {
	fmt.Fprintf(w, "define i32 @ugo_%s_init() {\n", file.Pkg.Name)

	for _, g := range file.Globals {
		var localName = "0"
		if g.Value != nil {
			localName = p.compileExpr(w, g.Value)
		}

		var varName string
		if _, obj := p.scope.Lookup(g.Name.Name); obj != nil {
			varName = obj.MangledName
		} else {
			panic(fmt.Sprintf("var %s undefined", g))
		}

		fmt.Fprintf(w, "\tstore i32 %s, i32* %s\n", localName, varName)
	}
	fmt.Fprintln(w, "\tret i32 0")
	fmt.Fprintln(w, "}")
}

func (p *Compiler) compileFile(w io.Writer, file *ast.File) {
	p.enterScope()
	defer p.leaveScope()

	for _, g := range file.Globals {
		var mangledName = fmt.Sprintf("@ugo_%s_%s", file.Pkg.Name, g.Name.Name)
		p.scope.Insert(&Object{
			Name:        g.Name.Name,
			MangledName: mangledName,
			Node:        g,
		})
		fmt.Fprintf(w, "%s = global i32 0\n", mangledName)
	}
	if len(file.Globals) != 0 {
		fmt.Fprintln(w)
	}
	for _, fn := range file.Funcs {
		p.compileFunc(w, file, fn)
	}

	p.genInit(w, file)
}

func (p *Compiler) compileFunc(w io.Writer, file *ast.File, fn *ast.Func) {
	p.enterScope()
	defer p.leaveScope()

	var mangledName = fmt.Sprintf("@ugo_%s_%s", file.Pkg.Name, fn.Name)

	p.scope.Insert(&Object{
		Name:        fn.Name,
		MangledName: mangledName,
		Node:        fn,
	})

	if fn.Body == nil {
		fmt.Fprintf(w, "declare i32 @ugo_%s_%s()\n", file.Pkg.Name, fn.Name)
		return
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "define i32 @ugo_%s_%s() {\n", file.Pkg.Name, fn.Name)
	p.compileStmt(w, fn.Body)

	fmt.Fprintln(w, "\tret i32 0")
	fmt.Fprintln(w, "}")
}

func (p *Compiler) compileStmt(w io.Writer, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.VarSpec:
		var localName = "0"
		if stmt.Value != nil {
			localName = p.compileExpr(w, stmt.Value)
		}

		var mangledName = fmt.Sprintf("%%local_%s.pos.%d", stmt.Name.Name, stmt.VarPos)
		p.scope.Insert(&Object{
			Name:        stmt.Name.Name,
			MangledName: mangledName,
			Node:        stmt,
		})

		fmt.Fprintf(w, "\t%s = alloca i32, align 4\n", mangledName)
		fmt.Fprintf(
			w, "\tstore i32 %s, i32* %s\n",
			localName, mangledName,
		)

	case *ast.AssignStmt:
		p.compileStmt_assign(w, stmt)
	case *ast.ForStmt:
		p.enterScope()
		defer p.leaveScope()

		forPos := fmt.Sprintf("%d", stmt.For)
		forInit := p.genLabelId("for.init.pos" + forPos)
		forCond := p.genLabelId("for.cond.pos" + forPos)
		forBody := p.genLabelId("for.body.pos" + forPos)
		forEnd := p.genLabelId("for.end.pos" + forPos)

		// for.init
		fmt.Fprintf(w, "\n%s:\n", forInit)
		if stmt.Init != nil {
			p.compileStmt(w, stmt.Init)
		}
		// br %forCond
		fmt.Fprintf(w, "\tbr label %%%s\n", forCond)

		fmt.Fprintf(w, "\n%s:\n", forCond)
		if stmt.Cond != nil {
			condValue := p.compileExpr(w, stmt.Cond)
			fmt.Fprintf(w, "\tbr i1 %s , label %%%s, label %%%s\n", condValue, forBody, forEnd)
		} else {
			fmt.Fprintf(w, "\tbr label %%%s\n", forBody)
		}

		// Body
		p.enterScope()
		fmt.Fprintf(w, "\n%s:\n", forBody)
		p.compileStmt(w, stmt.Body)
		p.leaveScope()

		// post
		if stmt.Post != nil {
			p.compileStmt(w, stmt.Post)
		}

		// br %forCond
		fmt.Fprintf(w, "\tbr label %%%s\n", forCond)

		// end
		fmt.Fprintf(w, "\n%s:\n", forEnd)

	case *ast.BlockStmt:
		p.enterScope()
		defer p.leaveScope()

		for _, x := range stmt.List {
			p.compileStmt(w, x)
		}
	case *ast.ExprStmt:
		p.compileExpr(w, stmt.X)

	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", stmt))
	}
}

func (p *Compiler) compileStmt_assign(w io.Writer, stmt *ast.AssignStmt) {
	var valueNameList = make([]string, len(stmt.Value))
	for i := range stmt.Target {
		valueNameList[i] = p.compileExpr(w, stmt.Value[i])
	}

	if stmt.Op == token.DEFINE {
		for _, target := range stmt.Target {
			if _, obj := p.scope.Lookup(target.Name); obj == nil {
				var mangledName = fmt.Sprintf("%%local_%s.pos.%d", target.Name, target.NamePos)
				p.scope.Insert(&Object{
					Name:        target.Name,
					MangledName: mangledName,
					Node:        target,
				})
				fmt.Fprintf(w, "\t%s = alloca i32, align 4\n", mangledName)
			}
		}
	}
	for i := range stmt.Target {
		var varName string
		if _, obj := p.scope.Lookup(stmt.Target[i].Name); obj != nil {
			varName = obj.MangledName
		} else {
			panic(fmt.Sprintf("var %s undefined", stmt.Target[0].Name))
		}

		fmt.Fprintf(
			w, "\tstore i32 %s, i32* %s\n",
			valueNameList[i], varName,
		)
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
		fmt.Fprintf(w, "\t%s = load i32, i32* %s, align 4\n",
			localName, varName,
		)
		return localName
	case *ast.Number:
		localName = p.genId()
		fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
			localName, "add", `0`, expr.Value,
		)
		return localName
	case *ast.BinaryExpr:
		localName = p.genId()
		switch expr.Op {
		case token.ADD:
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "add", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.SUB:
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "sub", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.MUL:
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "mul", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.DIV:
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "div", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName

		// https://llvm.org/docs/LangRef.html#icmp-instruction

		case token.EQL: // ==
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp eq", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.NEQ: // !=
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp ne", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.LSS: // <
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp slt", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.LEQ: // <=
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sle", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.GTR: // >
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sgt", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		case token.GEQ: // >=
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "icmp sge", p.compileExpr(w, expr.X), p.compileExpr(w, expr.Y),
			)
			return localName
		default:
			panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
		}
	case *ast.UnaryExpr:
		if expr.Op == token.SUB {
			localName = p.genId()
			fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n",
				localName, "sub", `0`, p.compileExpr(w, expr.X),
			)
			return localName
		}
		return p.compileExpr(w, expr.X)
	case *ast.ParenExpr:
		return p.compileExpr(w, expr.X)
	case *ast.CallExpr:
		var fnName string
		if _, obj := p.scope.Lookup(expr.FuncName.Name); obj != nil {
			fnName = obj.MangledName
		} else {
			panic(fmt.Sprintf("func %s undefined", expr.FuncName.Name))
		}

		localName = p.genId()
		fmt.Fprintf(w, "\t%s = call i32(i32) %s(i32 %v)\n",
			localName, fnName, p.compileExpr(w, expr.Args[0]),
		)
		return localName

	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
	}
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
