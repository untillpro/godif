package godif

// Provide implementation for pFunc in given context
func Provide(pFunc interface{}, impl interface{}) {
}

// Require implementation of function in given context
func Require(pFunc interface{}) {
}

// ResolveAll resolves all dependencies for RootCD and its sub-contexts
func ResolveAll() []error {
	return []error{}
}

// ProvideMapValue registers map value. pData must have a Key field
func ProvideMapValue(pMap interface{}, pData interface{}) {
}
