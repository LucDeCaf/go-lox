package lox

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt)
	VisitPrintStmt(*PrintStmt)
}

type ExpressionStmt struct {
	expression Expr
}

type PrintStmt struct {
	expression Expr
}

// TODO errors
func (s *ExpressionStmt) Accept(v StmtVisitor) error {
	v.VisitExpressionStmt(s)
	return nil
}

// TODO errors
func (s *PrintStmt) Accept(v StmtVisitor) error {
	v.VisitPrintStmt(s)
	return nil
}
