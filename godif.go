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

var required []*srcElem
var provided map[interface{}][]*srcElem
var keyValues map[interface{}]map[interface{}][]*srcElem
var sliceElements map[interface{}][]*srcElem

func init() {
	createVars()
}

func createVars() {
	provided = make(map[interface{}][]*srcElem)
	keyValues = make(map[interface{}]map[interface{}][]*srcElem)
	sliceElements = make(map[interface{}][]*srcElem)
}

// Reset clears all assignations
func Reset() {
	createVars()
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
		required = make([]*srcElem, 0)
	}
}

// ProvideSliceElement s.e.
func ProvideSliceElement(pointerToSlice interface{}, element interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if sliceElements[pointerToSlice] == nil {
		sliceElements[pointerToSlice] = make([]*srcElem, 0)
	}
	sliceElements[pointerToSlice] = append(sliceElements[pointerToSlice], &srcElem{file, line, element})
}

// ProvideKeyValue s.e.
func ProvideKeyValue(pointerToMap interface{}, key interface{}, value interface{}) {
	//requireEx(pMap, 2)
	_, file, line, _ := runtime.Caller(1)
	if keyValues[pointerToMap] == nil {
		keyValues[pointerToMap] = make(map[interface{}][]*srcElem)
	}
	keyValues[pointerToMap][key] = append(keyValues[pointerToMap][key], &srcElem{file, line, value})
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[ref] = append(provided[ref], &srcElem{file, line, funcImplementation})
}

// Require registers dep
func Require(toInject interface{}) {
	requireEx(toInject, 2)
}

// Require registers dep
func requireEx(toInject interface{}, callerStackOffset int) {
	_, file, line, _ := runtime.Caller(callerStackOffset)
	required = append(required, &srcElem{file, line, toInject})
}

// ResolveAll all deps
func ResolveAll() Errors {
	errs := getErrors()
	if len(errs) > 0 {
		return errs
	}

	for _, reqVar := range required {
		impls := provided[reqVar.elem]
		reqValue := reflect.ValueOf(reqVar.elem).Elem()

		reqType := reflect.TypeOf(reqVar.elem).Elem()
		reqKind := reqType.Kind()
		switch reqKind {
		case reflect.Func:
			reqValue.Set(reflect.ValueOf(impls[0].elem))
		case reflect.Map:
			reqValue.Set(reflect.ValueOf(impls[0].elem))
			kvToAppend := keyValues[reqVar.elem]
			for k, v := range kvToAppend {
				reqMapValueType := reqType.Elem()
				reqMapValueKind := reqMapValueType.Kind()
				provValue := reflect.ValueOf(v[0].elem)
				provValueKind := provValue.Kind()
				reqValue := reflect.ValueOf(reqVar.elem).Elem()
				keyValue := reflect.ValueOf(k)

				if !isSlice(provValueKind) && isSlice(reqMapValueKind) {
					provValue = reflect.New(reflect.SliceOf(reflect.TypeOf(v[0].elem))).Elem()
					for _, elementToAppend := range v {
						provElementValue := reflect.ValueOf(elementToAppend.elem)
						provValue.Set(reflect.Append(provValue, provElementValue))
					}
				}
				reqValue.SetMapIndex(keyValue, provValue)
			}
		case reflect.Slice, reflect.Array:
			sliceToAppend := sliceElements[reqVar.elem]
			for _, elementToAppend := range sliceToAppend {
				dataValue := reflect.ValueOf(elementToAppend.elem)
				reqValue := reflect.ValueOf(reqVar.elem).Elem()
				reqValue.Set(reflect.Append(reqValue, dataValue))
			}
		}
	}

	return nil
}

func isSlice(kind reflect.Kind) bool {
	return kind == reflect.Array || kind == reflect.Slice
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

		for _, v := range keyValues[req.elem] {
			reqMapValueType := reqType.Elem()
			reqMapValueKind := reqMapValueType.Kind()
			if isSlice(reqMapValueKind) {
				reqMapValueSliceElementType := reqMapValueType.Elem()
				for _, provElement := range v {
					provType := reflect.TypeOf(provElement.elem)
					provKind := provType.Kind()
					if isSlice(provKind) {
						if len(v) > 1 {
							errs = append(errs, &EMultipleValues{req, v})
							break
						}
						provType = provType.Elem()
						provKind = provType.Kind()
					}
					if !provType.AssignableTo(reqMapValueSliceElementType) {
						errs = append(errs, &EIncompatibleTypes{req, provElement})
					}
				}
			} else {
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

		for _, v := range sliceElements[req.elem] {
			vType := reflect.TypeOf(v.elem)
			if !vType.AssignableTo(reqType.Elem()) {
				errs = append(errs, &EIncompatibleTypes{req, v})
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
