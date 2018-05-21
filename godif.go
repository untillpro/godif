package godif

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

var registered map[reflect.Type][]interface{}
var required []interface{}

// to ContainerDeclaration
type request struct {
	pkgName string
	typ     reflect.Type
	impl    interface{}
}

func init() {
	registered = make(map[reflect.Type][]interface{})
}

// RegisterImpl register implementation
func RegisterImpl(funcImplementation interface{}) {
	RegisterImplByType(reflect.TypeOf(&funcImplementation), funcImplementation)
}

// RegisterImplByType registers implementation by type
func RegisterImplByType(typ reflect.Type, funcImplementation interface{}) {
	typ.Name()
	registered[typ] = append(registered[typ], funcImplementation)
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
	packagesProv, packagesReq := make([]string, 0), make([]string, 0)
	for _, reqVar := range required {
		test1 := reqVar.(interface{})
		fmt.Printf(reflect.TypeOf(test1).String())
		reqType := reflect.TypeOf(&reqVar)
		//fmt.Printf(reqType.String())
		impl := registered[reqType]
		fmt.Println(impl[0])
		fmt.Println(registered)
		if nil == impl {
			impl = registered[reflect.TypeOf(test1)]
			if nil == impl {
				return errors.New("required " + reqType.String() + " not registered")
			}
			// unresolved dependencies
			return errors.New("required " + reqType.String() + " not registered")
		}

		if len(impl) > 1 {
			// multiple implementations
			return fmt.Errorf("%s registered %d times", reqType.String(), len(impl))
		}

		pkgReq := reqType.PkgPath()
		pkgProv := reflect.TypeOf(impl[0]).PkgPath()

		packagesReq = append(packagesReq, pkgReq)
		packagesProv = append(packagesProv, pkgProv)

		v := reflect.ValueOf(reqVar).Elem()
		fmt.Println(v)
		v.Set(reflect.ValueOf(impl[0]))
	}
	return nil
}
