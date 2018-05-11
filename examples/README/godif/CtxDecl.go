package godif

import "context"

// CtxDecl -arations
type CtxDecl struct {
}

// Provide implementation for pFunc in given context
func (cd *CtxDecl) Provide(pFunc interface{}, impl interface{}) {
}

// Require implementation of function in given context
func (cd *CtxDecl) Require(pFunc interface{}) {
}

// ProvideMain function in given context
func (cd *CtxDecl) ProvideMain(pMain func(ctx context.Context)) {
}

// CreateCtxInst s.e.
func (cd *CtxDecl) CreateCtxInst(ctx context.Context) *CtxInst {
	return nil
}
