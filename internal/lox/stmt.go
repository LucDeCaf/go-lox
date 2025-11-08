package lox

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) error
	VisitPrintStmt(*PrintStmt) error
	VisitVarStmt(*VarStmt) error
}

type ExpressionStmt struct {
	expression Expr
}

type PrintStmt struct {
	expression Expr
}

type VarStmt struct {
	name  *Token
	value Expr
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
