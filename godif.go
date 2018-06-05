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
	errs := getErrors()
	if len(errs) > 0 {
		return errs
	}

	for _, reqVar := range required {
		v := reflect.ValueOf(reqVar.elem).Elem()
		impl := provided[reqVar.elem]
		v.Set(reflect.ValueOf(impl[0].elem))
	}
	return nil
}

func getErrors() Errors {
	var errs Errors
	for _, req := range required {

		kind := reflect.ValueOf(req.elem).Kind()

		if kind != reflect.Ptr {
			errs = append(errs, &ENonAssignableRequirement{req})
			continue
		}

		impls := provided[req.elem]

		if nil == impls {
			errs = append(errs, &EImplementationNotProvided{req})
		}

		if len(impls) > 1 {
			errs = append(errs, &EMultipleImplementations{req, impls})
		}

		v := reflect.ValueOf(req.elem).Elem()
		if !v.CanSet() {
			errs = append(errs, &ENonAssignableRequirement{req})
		}

		if len(impls) == 1 {
			reqType := reflect.TypeOf(req.elem).Elem()
			implType := reflect.TypeOf(impls[0].elem)
			if !implType.AssignableTo(reqType) {
				errs = append(errs, &EIncompatibleTypes{req, impls[0]})
			}
		}
	}
	return errs
}
