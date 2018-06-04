/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"reflect"
	"runtime"
)

type srcElem struct {
	file string
	line int
	elem interface{}
}

var provided map[reflect.Type][]srcElem
var required []srcElem

func init() {
	provided = make(map[reflect.Type][]srcElem)
}

// Reset clears all assignations
func Reset() {
	provided = make(map[reflect.Type][]srcElem)
	if required != nil {
		for _, r := range required {
			v := reflect.ValueOf(r.elem).Elem()
			v.Set(reflect.Zero(v.Type()))
		}
		required = make([]srcElem, 0)
	}
}

// ProvideByImpl registers implementation of funcImplementation type
func ProvideByImpl(funcImplementation interface{}) {
	ProvideByType(reflect.TypeOf(funcImplementation), funcImplementation)
}

// ProvideByType registers implementation by type
func ProvideByType(typ reflect.Type, funcImplementation interface{}) {
	provideByType(typ, funcImplementation)
}

func provideByType(typ reflect.Type, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(2)
	provided[typ] = append(provided[typ], srcElem{file, line, funcImplementation})
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	provideByType(reflect.TypeOf(ref).Elem(), funcImplementation)
}

// Require registers dep
func Require(pFunc interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, srcElem{file, line, pFunc})
}

// ResolveAll all deps
func ResolveAll() Errors {
	var errs Errors
	for _, reqVar := range required {
		reqType := reflect.TypeOf(reqVar.elem).Elem()
		impl := provided[reqType]

		if nil == impl {
			errs = append(errs, &EImplementationNotProvided{reqType})
		}

		if len(impl) > 1 {
			errs = append(errs, &EMultipleImplementations{reflect.TypeOf(reqVar.elem), impl})
		}

		v := reflect.ValueOf(reqVar.elem).Elem()

		if !v.CanSet() {
			errs = append(errs, &ENonAssignableRequirement{reqType, reqVar})
		}

		if len(errs) == 0 {
			v.Set(reflect.ValueOf(impl[0].elem))
		}
	}
	return errs
}
