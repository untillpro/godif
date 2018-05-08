package godif

// CtxInst -ance, context with data
type CtxInst struct {
}

// RegisterData in given CtxInst
func (ci *CtxInst) RegisterData(cd *ICtxData) {

}

// Init all registered CtxData, error is empty if panic does not happen
func (ci *CtxInst) Init(cd *ICtxData) []error {
	return []error{}
}
