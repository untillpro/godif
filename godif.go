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
	"strings"

	"github.com/untillpro/gochips/errs"
)

type src struct {
	file string
	line int
}

type srcElem struct {
	*src
	elem interface{}
}

type srcPkgElem struct {
	*srcElem
	pkg string
}

var (
	required        []*srcElem
	provided        map[interface{}][]*srcPkgElem
	keyValues       map[interface{}]map[interface{}][]*srcElem
	sliceElements   map[interface{}][]*srcElem
	resolveSrc      *src
	unhashableProvs []*src
)

func init() {
	createVars()
}

func createVars() {
	provided = make(map[interface{}][]*srcPkgElem)
	keyValues = make(map[interface{}]map[interface{}][]*srcElem)
	sliceElements = make(map[interface{}][]*srcElem)
}

func newSrcElem(file string, line int, elem interface{}) *srcElem {
	return &srcElem{&src{file, line}, elem}
}

func newSrcPkgElem(file string, line int, pkg string, elem interface{}) *srcPkgElem {
	return &srcPkgElem{newSrcElem(file, line, elem), pkg}
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
	for p := range provided {
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
	unhashableProvs = []*src{}
	createVars()
}

// ProvideSliceElement s.e.
func ProvideSliceElement(pointerToSlice interface{}, element interface{}) {
	_, file, line, _ := runtime.Caller(1)
	srcElement := newSrcElem(file, line, element)
	if isHashable(pointerToSlice) {
		sliceElements[pointerToSlice] = append(sliceElements[pointerToSlice], srcElement)
	} else {
		unhashableProvs = append(unhashableProvs, srcElement.src)
	}
}

// ProvideKeyValue s.e.
func ProvideKeyValue(pointerToMap interface{}, key interface{}, value interface{}) {
	_, file, line, _ := runtime.Caller(1)
	srcElement := newSrcElem(file, line, value)
	if isHashable(pointerToMap) {
		if keyValues[pointerToMap] == nil {
			keyValues[pointerToMap] = make(map[interface{}][]*srcElem)
		}
		keyValues[pointerToMap][key] = append(keyValues[pointerToMap][key], srcElement)
	} else {
		unhashableProvs = append(unhashableProvs, srcElement.src)
	}
}

// Provide registers implementation of ref type
func Provide(ref interface{}, funcImplementation interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	pkgName := nameFull[:strings.LastIndex(nameFull, ".")]
	srcElem := newSrcPkgElem(file, line, pkgName, funcImplementation)
	if isHashable(ref) {
		provided[ref] = append(provided[ref], srcElem)
	} else {
		unhashableProvs = append(unhashableProvs, srcElem.src)
	}
}

// Require registers dep
func Require(toInject interface{}) {
	_, file, line, _ := runtime.Caller(1)
	required = append(required, newSrcElem(file, line, toInject))
}

// ResolveAll all deps
func ResolveAll() errs.Errors {
	if errs := validate(); errs != nil {
		return errs
	}

	for target, provVar := range provided {
		if !targetRequired(target) {
			continue
		}
		if targetValue := reflect.ValueOf(target).Elem(); targetValue.IsNil() {
			targetValue.Set(reflect.ValueOf(provVar[0].elem))
		}
	}

	for targetMap, kvToAppend := range keyValues {
		targetMapType := reflect.TypeOf(targetMap).Elem()
		tragetMapValueType := targetMapType.Elem()
		tragetMapValueKind := tragetMapValueType.Kind()
		targetMapValue := reflect.ValueOf(targetMap).Elem()
		for k, v := range kvToAppend {
			keyValue := reflect.ValueOf(k)
			var toAppendValue reflect.Value
			if isSlice(tragetMapValueKind) {
				existingSlice := targetMapValue.MapIndex(keyValue)
				newSlice := reflect.New(reflect.SliceOf(tragetMapValueType.Elem())).Elem()
				if existingSlice.IsValid() {
					for i := 0; i < existingSlice.Len(); i++ {
						existingElement := existingSlice.Index(i)
						newSlice.Set(reflect.Append(newSlice, existingElement))
					}
				}
				for _, elementToAppend := range v {
					elementToAppendValue := reflect.ValueOf(elementToAppend.elem)
					elementToAppendKind := elementToAppendValue.Kind()
					if isSlice(elementToAppendKind) {
						for i := 0; i < elementToAppendValue.Len(); i++ {
							newSlice.Set(reflect.Append(newSlice, elementToAppendValue.Index(i)))
						}
					} else {
						newSlice.Set(reflect.Append(newSlice, elementToAppendValue))
					}
				}
				toAppendValue = newSlice
			} else {
				toAppendValue = reflect.ValueOf(v[0].elem)
			}
			targetMapValue.SetMapIndex(keyValue, toAppendValue)
		}
	}

	for targetSlice, elementsToAppend := range sliceElements {
		targateSliceValue := reflect.ValueOf(targetSlice).Elem()
		for _, elementToAppend := range elementsToAppend {
			elementValue := reflect.ValueOf(elementToAppend.elem)
			elementKind := elementValue.Kind()
			if isSlice(elementKind) {
				for i := 0; i < elementValue.Len(); i++ {
					targateSliceValue.Set(reflect.Append(targateSliceValue, elementValue.Index(i)))
				}
			} else {
				targateSliceValue.Set(reflect.Append(targateSliceValue, elementValue))
			}
		}
	}

	_, file, line, _ := runtime.Caller(1)
	resolveSrc = &src{file, line}

	return nil
}

func targetRequired(target interface{}) bool {
	for _, elem := range required {
		if elem.elem == target {
			return true
		}
	}
	return false
}

func isSlice(kind reflect.Kind) bool {
	return kind == reflect.Array || kind == reflect.Slice
}

func isHashable(intf interface{}) bool {
	k := reflect.TypeOf(intf).Kind()
	return k < reflect.Array || k == reflect.Ptr || k == reflect.UnsafePointer
}

func validate() (errs errs.Errors) {
	if resolveSrc != nil {
		return errs.AddE(&EAlreadyResolved{resolveSrc})
	}

	requiredPackages := make(map[string]bool)

	if len(unhashableProvs) > 0 {
		for _, unhashableProvsSrc := range unhashableProvs {
			errs.AddE(&EProvisionForNonAssignable{unhashableProvsSrc})
		}
		return errs
	}

	for _, req := range required {

		v := reflect.ValueOf(req.elem)

		if v.Kind() != reflect.Ptr || !v.Elem().CanSet() {
			errs.AddE(&ENonAssignableRequirement{req})
			if !v.CanSet() {
				return errs // req.elem is unhashable here
			}
		}

		impls := provided[req.elem]

		if nil == impls {
			errs.AddE(&EImplementationNotProvided{req, nil})
		}

		if len(impls) > 1 {
			errs.AddE(&EMultipleFuncImplementations{req, impls})
		}

		reqType := reflect.TypeOf(req.elem).Elem()

		for _, impl := range impls {
			requiredPackages[impl.pkg] = true
			implType := reflect.TypeOf(impl.elem)
			if !implType.AssignableTo(reqType) {
				errs.AddE(&EIncompatibleTypesFunc{req, impl})
			}
		}
	}

	for targetMap, kvToAppend := range keyValues {
		targetMapType := reflect.TypeOf(targetMap).Elem()
		targetMapValue := reflect.ValueOf(targetMap).Elem()
		targetMapKeyType := targetMapType.Key()
		impl := provided[targetMap]
		if targetMapValue.IsNil() {
			if impl == nil {
				keys := reflect.ValueOf(kvToAppend).MapKeys()
				errs.AddE(&EImplementationNotProvided{kvToAppend[keys[0].Interface()][0], targetMap})
				continue
			}
		} else {
			if impl != nil {
				errs.AddE(&EImplementationProvidedForNonNil{impl[0]})
				continue
			}
		}
		targetMapValueType := targetMapType.Elem()
		targetMapValueKind := targetMapValueType.Kind()
		for k, v := range kvToAppend {
			if isSlice(targetMapValueKind) {
				reqMapValueSliceElementType := targetMapValueType.Elem()
				for _, provElement := range v {
					provType := reflect.TypeOf(provElement.elem)
					provKind := provType.Kind()
					if isSlice(provKind) {
						provType = provType.Elem()
					}
					if !provType.AssignableTo(reqMapValueSliceElementType) {
						errs.AddE(&EIncompatibleTypesStorageValue{targetMapType, provElement})
					}
				}
			} else {
				if len(v) > 1 {
					errs.AddE(&EMultipleValues{v})
				} else {
					vType := reflect.TypeOf(v[0].elem)
					if !vType.AssignableTo(targetMapValueType) {
						errs.AddE(&EIncompatibleTypesStorageValue{targetMapType, v[0]})
					}
					kType := reflect.TypeOf(k)
					if !kType.AssignableTo(targetMapKeyType) {
						errs.AddE(&EIncompatibleTypesStorageKey{targetMapType, newSrcElem(v[0].file, v[0].line, k)})
					}
				}
			}
		}
	}

	for targetSlice, elementsToAppend := range sliceElements {
		targetSliceType := reflect.TypeOf(targetSlice).Elem()
		targetSliceValue := reflect.ValueOf(targetSlice).Elem()
		impl := provided[targetSlice]
		if targetSliceValue.IsNil() && impl == nil {
			errs.AddE(&EImplementationNotProvided{elementsToAppend[0], targetSlice})
			continue
		}
		for _, v := range elementsToAppend {
			vType := reflect.TypeOf(v.elem)
			vKind := vType.Kind()
			if isSlice(vKind) {
				vType = vType.Elem()
			}
			if !vType.AssignableTo(targetSliceType.Elem()) {
				errs.AddE(&EIncompatibleTypesStorageValue{targetSliceType, v})
			}
		}
	}

	pkgNotUsedErrorsAppended := make(map[string]bool)

	for provVar, provSrcs := range provided {
		provKind := reflect.TypeOf(provVar).Elem().Kind()
		if provKind != reflect.Func && len(provSrcs) > 1 {
			errs.AddE(&EMultipleStorageImplementations{provSrcs})
			continue
		}
		provType := reflect.TypeOf(provSrcs[0].elem)
		targetType := reflect.TypeOf(provVar).Elem()
		targetKind := targetType.Kind()

		switch targetKind {
		case reflect.Func:
			if _, required := requiredPackages[provSrcs[0].pkg]; !required {
				if !pkgNotUsedErrorsAppended[provSrcs[0].pkg] {
					errs.AddE(&EPackageNotUsed{provSrcs[0].pkg})
					pkgNotUsedErrorsAppended[provSrcs[0].pkg] = true
				}
			}
		case reflect.Array, reflect.Slice, reflect.Map:
			if isSlice(targetKind) {
				targetSliceValue := reflect.ValueOf(provVar).Elem()
				if !targetSliceValue.IsNil() {
					errs.AddE(&EImplementationProvidedForNonNil{provSrcs[0]})
				}
			}
			if !provType.AssignableTo(targetType) {
				errs.AddE(&EIncompatibleTypesStorageImpl{targetType, provSrcs[0].srcElem})
			}
		}
	}

	return errs
}
