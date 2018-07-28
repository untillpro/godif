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

var required []srcElem
var provided map[interface{}][]srcElem
var providedMapValues map[interface{}]map[interface{}][]srcElem

func init() {
	provided = make(map[interface{}][]srcElem)
	providedMapValues = make(map[interface{}]map[interface{}][]srcElem)
}

func reset() {
	provided = make(map[interface{}][]srcElem)
	providedMapValues = make(map[interface{}]map[interface{}][]srcElem)
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

func provideMapValue(pMap interface{}, key interface{}, data interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if providedMapValues[pMap] == nil {
		providedMapValues[pMap] = make(map[interface{}][]srcElem)
		fmt.Println(providedMapValues[pMap])
	}
	providedMapValues[pMap][key] = append(providedMapValues[pMap][key], srcElem{file, line, data})
}

func provide(ref interface{}, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[ref] = append(provided[ref], srcElem{file, line, funcImplementation})
}

func require(toInject interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, srcElem{file, line, toInject})
}

func resolveAll() Errors {
	errs := getErrors()
	if len(errs) > 0 {
		return errs
	}

	for _, reqVar := range required {
		impls := provided[reqVar.elem]
		reqValue := reflect.ValueOf(reqVar.elem).Elem()
		reqValue.Set(reflect.ValueOf(impls[0].elem))
	}
	for _, reqVar := range required {
		mapToAppend := providedMapValues[reqVar.elem]
		for k, v := range mapToAppend {
			dataValue := reflect.ValueOf(v[0].elem)
			reqValue := reflect.ValueOf(reqVar.elem).Elem()
			keyValue := reflect.ValueOf(k)
			reqValue.SetMapIndex(keyValue, dataValue)
		}
	}

	return nil
}

func getErrors() Errors {
	var errs Errors
	for _, req := range required {

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

		reqType := reflect.TypeOf(req.elem).Elem()

		for _, impl := range impls {
			implType := reflect.TypeOf(impl.elem)
			if !implType.AssignableTo(reqType) {
				errs = append(errs, &EIncompatibleTypes{req, impls[0]})
			}
		}

		for _, v := range providedMapValues[req.elem] {
			if len(v) > 1 {
				errs = append(errs, &EMultipleValues{req, v})
			} else {
				vType := reflect.TypeOf(v[0].elem)
				if !vType.AssignableTo(reqType.Elem()) {
					errs = append(errs, &EIncompatibleTypes{req, v[0]})
				}
			}
		}
	}
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
