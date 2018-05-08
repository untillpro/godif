package godif

// CtxDecl -arations
type CtxDecl struct {
}

// Provide implementation for pFunc in given context
func (cd *CtxDecl) Provide(pFunc interface{}, impl interface{}) {
}

// Require implementation of function in given context
func (cd *CtxDecl) Require(pFunc interface{}) {
}

// RegisterCtxMain function in given context
func (cd *CtxDecl) RegisterCtxMain(pInitFunc interface{}) {
}

// CreateCtxInst s.e.
func (cd *CtxDecl) CreateCtxInst() *CtxInst {
	return nil
}
