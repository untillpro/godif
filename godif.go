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

var provided map[interface{}][]srcElem
var required []srcElem

func init() {
	provided = make(map[interface{}][]srcElem)
}

// Reset clears all assignations
func Reset() {
	provided = make(map[interface{}][]srcElem)
	if required != nil {
		for _, r := range required {
			v := reflect.ValueOf(r.elem)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
				if v.CanSet() {
					v.Set(reflect.Zero(v.Type()))
				}
			}
		}
		required = make([]srcElem, 0)
	}
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[ref] = append(provided[ref], srcElem{file, line, funcImplementation})
}

// Require registers dep
func Require(toInject interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, srcElem{file, line, toInject})
}

// ResolveAll all deps
func ResolveAll() Errors {
	var errs Errors
	for _, reqVar := range required {

		kind := reflect.ValueOf(reqVar.elem).Kind()

		if kind != reflect.Ptr {
			errs = append(errs, &ENonAssignableRequirement{reqVar})
			continue
		}

		v := reflect.ValueOf(reqVar.elem).Elem()

		impl := provided[reqVar.elem]

		if nil == impl {
			errs = append(errs, &EImplementationNotProvided{reqVar})
		}

		if len(impl) > 1 {
			errs = append(errs, &EMultipleImplementations{reflect.TypeOf(reqVar.elem), impl})
		}

		if !v.CanSet() {
			errs = append(errs, &ENonAssignableRequirement{reqVar})
		}

		if len(errs) == 0 {
			v.Set(reflect.ValueOf(impl[0].elem))
		}
	}
	return errs
}
