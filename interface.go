package godif

// Reset clears all assignations
var Reset func()

// ProvideMapValue s.e.
var ProvideMapValue func(pMap interface{}, key interface{}, data interface{})

// ProvideSliceValue s.e.
var ProvideSliceValue func(pSlice interface{}, data interface{})

// Provide registers implementation of ref type
var Provide func(ref interface{}, funcImplementation interface{})

// Require registers dep
var Require func(toInject interface{})

// ResolveAll all deps
var ResolveAll func() Errors

func init() {
	Reset = reset
	ProvideMapValue = provideMapValue
	ProvideSliceValue = provideSliceValue
	Provide = provide
	Require = require
	ResolveAll = resolveAll
}
