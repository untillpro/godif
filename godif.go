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
	file   string
	line   int
	elem   interface{}
	isData bool
}

var required []srcElem
var provided map[interface{}][]srcElem

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

// ProvideMapValue registers map value. pData must have a Key field
func ProvideMapValue(pMap interface{}, pData interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[pMap] = append(provided[pMap], srcElem{file, line, pData, true})
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[ref] = append(provided[ref], srcElem{file, line, funcImplementation, false})
}

// Require registers dep
func Require(toInject interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, srcElem{file, line, toInject, false})
}

// ResolveAll all deps
func ResolveAll() Errors {
	errs := getErrors()
	if len(errs) > 0 {
		return errs
	}

	for _, reqVar := range required {
		impls := provided[reqVar.elem]
		for _, impl := range impls {
			if impl.isData {
				data := reflect.ValueOf(impl.elem)
				reqValue := reflect.ValueOf(reqVar.elem).Elem()
				if data.Kind() == reflect.Ptr {
					reqValue.SetMapIndex(data.Elem().FieldByName("Key"), data)
				} else if data.Kind() == reflect.Struct {
					reqValue.SetMapIndex(data.FieldByName("Key"), data)
				}
			} else {
				reqValue := reflect.ValueOf(reqVar.elem).Elem()
				reqValue.Set(reflect.ValueOf(impl.elem))
			}
		}
	}

	return nil
}

func getErrors() Errors {
	var errs Errors
	for _, req := range required {

		//kind := reflect.ValueOf(req.elem).Kind()

		// if kind != reflect.Ptr {
		// 	errs = append(errs, &ENonAssignableRequirement{req})
		// 	continue
		// }

		impls := provided[req.elem]

		if nil == impls {
			errs = append(errs, &EImplementationNotProvided{req})
		}

		if len(impls) > 2 || (len(impls) == 2 && impls[0].isData == impls[1].isData) {
			errs = append(errs, &EMultipleImplementations{req, impls})
		}

		v := reflect.ValueOf(req.elem).Elem()
		if !v.CanSet() {
			errs = append(errs, &ENonAssignableRequirement{req})
		}

		for _, impl := range impls {
			reqType := reflect.TypeOf(req.elem).Elem()
			fmt.Println(reqType)
			implType := reflect.TypeOf(impl.elem)
			fmt.Println(implType)
			fmt.Println(reqType.Kind())
			if impl.isData {
				if !implType.AssignableTo(reqType.Elem()) {
					errs = append(errs, &EIncompatibleTypes{req, impl})
				}
			} else {
				if !implType.AssignableTo(reqType) {
					errs = append(errs, &EIncompatibleTypes{req, impls[0]})
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
