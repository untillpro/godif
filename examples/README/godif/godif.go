package godif


// RootCD is a root context declaration
var RootCD = &CtxDecl{}

// ResolveAll resolves all dependencies for RootCD and its sub-contexts
func ResolveAll() []error {
	return []error{}
}
