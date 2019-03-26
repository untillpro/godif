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

type src struct {
	file string
	line int
}

type srcElem struct {
	*src
	elem interface{}
}

var required []*srcElem
var provided map[interface{}][]*srcElem
var keyValues map[interface{}]map[interface{}][]*srcElem
var sliceElements map[interface{}][]*srcElem
var resolveSrc *src = nil

func init() {
	createVars()
}

func createVars() {
	provided = make(map[interface{}][]*srcElem)
	keyValues = make(map[interface{}]map[interface{}][]*srcElem)
	sliceElements = make(map[interface{}][]*srcElem)
}

func newSrcElem(file string, line int, elem interface{}) *srcElem {
	return &srcElem{&src{file, line}, elem}
}

// Reset clears all assignations
func Reset() {
	for _, r := range required {
		v := reflect.ValueOf(r.elem)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			if v.CanSet() {
				v.Set(reflect.Zero(v.Type()))
			}
		}
	}
	for p, _ := range provided {
		v := reflect.ValueOf(p)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			if v.CanSet() {
				v.Set(reflect.Zero(v.Type()))
			}
		}
	}
	required = make([]*srcElem, 0)
	resolveSrc = nil
	createVars()
}

// ProvideSliceElement s.e.
func ProvideSliceElement(pointerToSlice interface{}, element interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if sliceElements[pointerToSlice] == nil {
		sliceElements[pointerToSlice] = make([]*srcElem, 0)
	}
	sliceElements[pointerToSlice] = append(sliceElements[pointerToSlice], newSrcElem(file, line, element))
}

// ProvideKeyValue s.e.
func ProvideKeyValue(pointerToMap interface{}, key interface{}, value interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if keyValues[pointerToMap] == nil {
		keyValues[pointerToMap] = make(map[interface{}][]*srcElem)
	}
	keyValues[pointerToMap][key] = append(keyValues[pointerToMap][key], newSrcElem(file, line, value))
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	_, file, line, _ := runtime.Caller(1)
	provided[ref] = append(provided[ref], newSrcElem(file, line, funcImplementation))
}

// Require registers dep
func Require(toInject interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, newSrcElem(file, line, toInject))
}

// ResolveAll all deps
func ResolveAll() Errors {
	errs := getErrors()
	if len(errs) > 0 {
		return errs
	}

	for target, provVar := range provided {
		targetValue := reflect.ValueOf(target).Elem()
		if targetValue.IsNil() {
			targetValue.Set(reflect.ValueOf(provVar[0].elem))
		}
	}

	for targetMap, kvToAppend := range keyValues {
		targetMapType := reflect.TypeOf(targetMap).Elem()
		tragetMapValueType := targetMapType.Elem()
		tragetMapValueKind := tragetMapValueType.Kind()
		targetMapValue := reflect.ValueOf(targetMap).Elem()
		for k, v := range kvToAppend {
			valueValue := reflect.ValueOf(v[0].elem) // "str" -> [1 v0, 2 v1]
			valueValueKind := valueValue.Kind()
			keyValue := reflect.ValueOf(k)
			if !isSlice(valueValueKind) && isSlice(tragetMapValueKind) {
				existingSlice := targetMapValue.MapIndex(keyValue)
				newSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(v[0].elem))).Elem()
				if existingSlice.IsValid() {
					for i:=0; i< existingSlice.Len(); i++ {
						existingElement := existingSlice.Index(i)
						newSlice.Set(reflect.Append(newSlice, existingElement))
					}
				} 
				for _, elementToAppend := range v {
					elementToAppendValue := reflect.ValueOf(elementToAppend.elem)
					newSlice.Set(reflect.Append(newSlice, elementToAppendValue))
				}
				valueValue = newSlice
			}
			targetMapValue.SetMapIndex(keyValue, valueValue)
		}
	}

	for targetSlice, elementsToAppend := range sliceElements {
		for _, elementToAppend := range elementsToAppend {
			elementValue := reflect.ValueOf(elementToAppend.elem)
			reqValue := reflect.ValueOf(targetSlice).Elem()
			reqValue.Set(reflect.Append(reqValue, elementValue))
		}
	}

	_, file, line, _ := runtime.Caller(1)
	resolveSrc = &src{file, line}

	return nil
}

func isSlice(kind reflect.Kind) bool {
	return kind == reflect.Array || kind == reflect.Slice
}

func getErrors() Errors {
	var errs Errors
	if resolveSrc != nil {
		return []error{&EAlreadyResolved{resolveSrc}}
	}
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
				errs = append(errs, &EIncompatibleTypesFunc{req, impls[0]})
			}
		}
	}

	for targetMap, kvToAppend := range keyValues {
		targetMapType := reflect.TypeOf(targetMap).Elem()
		targetValue := reflect.ValueOf(targetMap).Elem()
		impl := provided[targetMap]
		if targetValue.IsNil() {
			if impl == nil {
				keys := reflect.ValueOf(kvToAppend).MapKeys()
				errs = append(errs, &EImplementationNotProvided{kvToAppend[keys[0].Interface()][0]})
				continue
			}
		} else {
			if impl != nil {
				errs = append(errs, &EImplementationProvidedForNonNil{impl[0]})
				continue
			}
		}
		tragetMapValueType := targetMapType.Elem()
		tragetMapValueKind := tragetMapValueType.Kind()
		for _, v := range kvToAppend {
			if isSlice(tragetMapValueKind) {
				reqMapValueSliceElementType := tragetMapValueType.Elem()
				for _, provElement := range v {
					provType := reflect.TypeOf(provElement.elem)
					provKind := provType.Kind()
					if isSlice(provKind) {
						if len(v) > 1 {
							errs = append(errs, &EMultipleValues{v})
							break
						}
						provType = provType.Elem()
						provKind = provType.Kind()
					}
					if !provType.AssignableTo(reqMapValueSliceElementType) {
						errs = append(errs, &EIncompatibleTypesSlice{targetMapType, provElement})
					}
				}
			} else {
				if len(v) > 1 {
					errs = append(errs, &EMultipleValues{v})
				} else {
					vType := reflect.TypeOf(v[0].elem)
					if !vType.AssignableTo(tragetMapValueType) {
						errs = append(errs, &EIncompatibleTypesSlice{tragetMapValueType, v[0]})
					}
				}
			}
		}
	}

	for targetSlice, elementsToAppend := range sliceElements {
		targetSliceType := reflect.TypeOf(targetSlice).Elem()
		targetSliceValue := reflect.ValueOf(targetSlice).Elem()
		impl := provided[targetSlice]
		if targetSliceValue.IsNil() {
			if impl == nil {
				errs = append(errs, &EImplementationNotProvided{elementsToAppend[0]})
				continue
			}
		} else {
			if impl != nil {
				errs = append(errs, &EImplementationProvidedForNonNil{impl[0]})
				continue
			}
		}
		for _, v := range elementsToAppend {
			vType := reflect.TypeOf(v.elem)
			if !vType.AssignableTo(targetSliceType.Elem()) {
				errs = append(errs, &EIncompatibleTypesSlice{targetSliceType, v})
			}
		}
	}

	// funcs only
	for provVar, provSrcs := range provided {
		if reflect.TypeOf(provVar).Elem().Kind() != reflect.Func {
			continue
		}
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
