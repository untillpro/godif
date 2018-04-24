package godif

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

var funcs map[reflect.Type]interface{}
var required []interface{}

func init() {
	funcs = make(map[reflect.Type]interface{})
}

// RegisterImpl register implementation
func RegisterImpl(funcImplementation interface{}) {
	RegisterImplByType(reflect.TypeOf(funcImplementation), funcImplementation)
}

// RegisterImplByType registers implementation by type
func RegisterImplByType(typ reflect.Type, funcImplementation interface{}) {
	funcs[typ] = funcImplementation
	funcImplType := reflect.TypeOf(funcImplementation)
	log.Println("Registered:", funcImplType, "pkg=", funcImplType.PkgPath())
}

// RegisterImplByVar registers implementation by nil var
func RegisterImplByVar(ref interface{}, funcImplementation interface{}) {
	RegisterImplByType(reflect.TypeOf(ref), funcImplementation)
}

// Require registers dep
func Require(pFunc interface{}) {
	required = append(required, pFunc)
	log.Println("Required:", reflect.TypeOf(pFunc), "pkg=", reflect.TypeOf(pFunc).PkgPath())
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fmt.Println("Called from", details.Name())
	}
}

// ResolveAll all deps
func ResolveAll() error {
	for _, pFunc := range required {
		t := reflect.TypeOf(pFunc).Elem()
		f := funcs[t]
		if nil == f {
			log.Panicln("required ", t, " not registered")
			return errors.New("required " + t.Name() + " not registered")
		}

		v := reflect.ValueOf(pFunc).Elem()
		v.Set(reflect.ValueOf(f))
	}
	return nil
}
