package godif

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

var registered map[reflect.Type][]interface{}
var required []interface{}

func init() {
	registered = make(map[reflect.Type][]interface{})
}

// Reset clears all assignations
func Reset() {
	registered = make(map[reflect.Type][]interface{})
	if required != nil {
		for _, r := range required {
			v := reflect.ValueOf(r).Elem()
			v.Set(reflect.Zero(v.Type()))
		}
		required = make([]interface{}, 0)
	} 
}

// ProvideByImpl register implementation of funcImplementation type
func ProvideByImpl(funcImplementation interface{}) {
	ProvideByType(reflect.TypeOf(funcImplementation), funcImplementation)
}

// ProvideByType registers implementation by type
func ProvideByType(typ reflect.Type, funcImplementation interface{}) {
	registered[typ] = append(registered[typ], funcImplementation)
	funcImplType := reflect.TypeOf(funcImplementation)
	log.Println("Registered:", funcImplType)
}

// Provide registers implementation of ref type 
func Provide(ref interface{}, funcImplementation interface{}) {
	ProvideByType(reflect.TypeOf(ref).Elem(), funcImplementation)
}

// Require registers dep
func Require(pFunc interface{}) {
	required = append(required, pFunc)
	log.Println("Required:", reflect.TypeOf(pFunc))
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fmt.Println("Called from", details.Name())
	}
}

// ResolveAll all deps
func ResolveAll() Errors {
	var errs Errors
	for _, reqVar := range required {
		reqType := reflect.TypeOf(reqVar).Elem()
		impl := registered[reqType]

		if nil == impl {
			errs = append(errs, &EImplementationNotProvided{reqType})
		}

		if len(impl) > 1 {
			errs = append(errs, &EMultipleImplementations{reflect.TypeOf(reqVar), len(impl)})
		}

		v := reflect.ValueOf(reqVar).Elem()

		if !v.CanSet() {
			errs = append(errs, &ENonAssignableRequirement{reqType})
		}

		if len(errs) == 0 {
			v.Set(reflect.ValueOf(impl[0]))
		}
	}
	return errs
}
