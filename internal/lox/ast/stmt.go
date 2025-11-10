package ast

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) error
	VisitPrintStmt(*PrintStmt) error
	VisitVarStmt(*VarStmt) error
}

type ExpressionStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name  *Token
	Value Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) error {
	return v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v StmtVisitor) error {
	return v.VisitPrintStmt(s)
}

func (s *VarStmt) Accept(v StmtVisitor) error {
	return v.VisitVarStmt(s)
}
