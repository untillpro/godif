package godif

// Reset clears all assignations
var Reset func()

// ProvideMapValue registers data which will be set on pMap map by "key" key on ResolveAll() call
var ProvideMapValue func(pMap interface{}, key interface{}, data interface{})

// Provide registers implementation of ref type
var Provide func(ref interface{}, funcImplementation interface{})

// Require registers dep
var Require func(toInject interface{})

// ResolveAll all deps
var ResolveAll func() Errors

func init() {
	Reset = reset
	ProvideMapValue = provideMapValue
	Provide = provide
	Require = require
	ResolveAll = resolveAll
}
