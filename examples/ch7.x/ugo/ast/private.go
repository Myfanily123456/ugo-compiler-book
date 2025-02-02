package ast

import "github.com/chai2010/ugo/token"

var (
	_ Node = Expr(nil)
	_ Node = Stmt(nil)

	_ Node = (*File)(nil)

	_ Node = (*PackageSpec)(nil)
	_ Node = (*ImportSpec)(nil)

	_ Stmt = (*ConstSpec)(nil)
	_ Stmt = (*VarSpec)(nil)
	_ Stmt = (*Func)(nil)

	_ Stmt = (*BlockStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)

	_ Expr = (*Ident)(nil)
	_ Expr = (*BasicLit)(nil)
	_ Expr = (*Number)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
	_ Expr = (*CallExpr)(nil)
	_ Expr = (*SelectorExpr)(nil)
)

func (p *File) Pos() token.Pos { return token.NoPos }
func (p *File) End() token.Pos { return token.NoPos }
func (p *File) node_type()     {}

func (p *PackageSpec) Pos() token.Pos { return token.NoPos }
func (p *PackageSpec) End() token.Pos { return token.NoPos }
func (p *PackageSpec) node_type()     {}

func (p *ImportSpec) Pos() token.Pos { return token.NoPos }
func (p *ImportSpec) End() token.Pos { return token.NoPos }
func (p *ImportSpec) node_type()     {}

func (p *ConstSpec) Pos() token.Pos { return token.NoPos }
func (p *ConstSpec) End() token.Pos { return token.NoPos }
func (p *ConstSpec) node_type()     {}

func (p *VarSpec) Pos() token.Pos { return token.NoPos }
func (p *VarSpec) End() token.Pos { return token.NoPos }
func (p *VarSpec) node_type()     {}

func (p *Func) Pos() token.Pos { return token.NoPos }
func (p *Func) End() token.Pos { return token.NoPos }
func (p *Func) node_type()     {}

func (p *BlockStmt) node_type()  {}
func (p *ExprStmt) node_type()   {}
func (p *AssignStmt) node_type() {}

func (p *Ident) node_type()        {}
func (p *BasicLit) node_type()     {}
func (p *Number) node_type()       {}
func (p *BinaryExpr) node_type()   {}
func (p *UnaryExpr) node_type()    {}
func (p *ParenExpr) node_type()    {}
func (p *CallExpr) node_type()     {}
func (p *SelectorExpr) node_type() {}

func (p *ConstSpec) stmt_type() {}
func (p *VarSpec) stmt_type()   {}
func (p *Func) stmt_type()      {}

func (p *BlockStmt) stmt_type()  {}
func (p *ExprStmt) stmt_type()   {}
func (p *AssignStmt) stmt_type() {}

func (p *Ident) expr_type()        {}
func (p *BasicLit) expr_type()     {}
func (p *Number) expr_type()       {}
func (p *BinaryExpr) expr_type()   {}
func (p *UnaryExpr) expr_type()    {}
func (p *ParenExpr) expr_type()    {}
func (p *CallExpr) expr_type()     {}
func (p *SelectorExpr) expr_type() {}

func (p *BlockStmt) Pos() token.Pos  { return token.NoPos }
func (p *ExprStmt) Pos() token.Pos   { return token.NoPos }
func (p *AssignStmt) Pos() token.Pos { return token.NoPos }

func (p *Ident) Pos() token.Pos        { return token.NoPos }
func (p *BasicLit) Pos() token.Pos     { return token.NoPos }
func (p *Number) Pos() token.Pos       { return token.NoPos }
func (p *BinaryExpr) Pos() token.Pos   { return token.NoPos }
func (p *UnaryExpr) Pos() token.Pos    { return token.NoPos }
func (p *ParenExpr) Pos() token.Pos    { return token.NoPos }
func (p *CallExpr) Pos() token.Pos     { return token.NoPos }
func (p *SelectorExpr) Pos() token.Pos { return token.NoPos }

func (p *BlockStmt) End() token.Pos  { return token.NoPos }
func (p *ExprStmt) End() token.Pos   { return token.NoPos }
func (p *AssignStmt) End() token.Pos { return token.NoPos }

func (p *Ident) End() token.Pos        { return token.NoPos }
func (p *BasicLit) End() token.Pos     { return token.NoPos }
func (p *Number) End() token.Pos       { return token.NoPos }
func (p *BinaryExpr) End() token.Pos   { return token.NoPos }
func (p *UnaryExpr) End() token.Pos    { return token.NoPos }
func (p *ParenExpr) End() token.Pos    { return token.NoPos }
func (p *CallExpr) End() token.Pos     { return token.NoPos }
func (p *SelectorExpr) End() token.Pos { return token.NoPos }
