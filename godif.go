/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"fmt"
	"reflect"
	"runtime"
)

type srcElem struct {
	file string
	line int
	elem interface{}
}

var provided map[interface{}][]srcElem
var providedMapValues map[interface{}][]srcElem
var required []srcElem

func init() {
	provided = make(map[interface{}][]srcElem)
	providedMapValues = make(map[interface{}][]srcElem)
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

// ProvideMapValue registers map value. pData must have a Key field
func ProvideMapValue(pMap interface{}, pData interface{}) {
	_, file, line, _ := runtime.Caller(1)
	providedMapValues[pMap] = append(providedMapValues[pMap], srcElem{file, line, pData})
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

	for ptrToVar, mapValues := range providedMapValues {
		v := reflect.ValueOf(mapValues[0].elem) // is struct bucketDef
		fmt.Println(v.Type())
		fmt.Println(v.Kind())
		fmt.Println(ptrToVar)
		m := reflect.ValueOf(ptrToVar).Elem()
		fmt.Println(m)
		m.SetMapIndex(v.Elem().FieldByName("Key"), v)
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
	if errsToAppend := getProvidedNotUsedErrors(provided); len(errsToAppend) > 0 {
		errs = append(errs, errsToAppend)
	}
	if errsToAppend := getProvidedNotUsedErrors(providedMapValues); len(errsToAppend) > 0 {
		errs = append(errs, errsToAppend)
	}
	return errs
}

func getProvidedNotUsedErrors(provided map[interface{}][]srcElem) Errors {
	var errs Errors
	for provVar, provSrcs := range provided {
		var found = false
		for _, reqVar := range required {
			if reqVar.elem == provVar {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, &EProvidedNotUsed{provSrcs[0]})
		}
	}
	return errs
}
