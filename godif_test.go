/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncBasic(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	Provide(&injectedFunc, f)

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, injectedFunc(3, 2))

	Reset()
	assert.Nil(t, injectedFunc)
}

func TestFuncErrorOnNoImplementation(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	errs := ResolveAll()

	if _, ok := errs[0].(*EImplementationNotProvided); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal(errs)
	}

	assert.Nil(t, injectedFunc)
}

func TestExplicitTypeInject(t *testing.T) {
	Reset()
	type Func1Type = func(x int, y int) int
	var inject Func1Type

	Require(&inject)
	Provide(&inject, f)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, inject(2, 3))
}

func TestErrorOnMultipleImplementations(t *testing.T) {
	Reset()
	var injectedFunc1 func(x int, y int) int

	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&injectedFunc1)
	_, implFileF, implLineF, _ := runtime.Caller(0)
	Provide(&injectedFunc1, f)
	_, implFileF2, implLineF2, _ := runtime.Caller(0)
	Provide(&injectedFunc1, f3)

	errs := ResolveAll()

	if e, ok := errs[0].(*EMultipleFuncImplementations); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, 2, len(e.provs))
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, implLineF+1, e.provs[0].line)
		assert.Equal(t, implFileF, e.provs[0].file)
		assert.Equal(t, implLineF2+1, e.provs[1].line)
		assert.Equal(t, implFileF2, e.provs[1].file)
	} else {
		t.Fatal(errs)
	}
	assert.Nil(t, injectedFunc1)
}

func TestMultipleErrorsOnResolve(t *testing.T) {
	Reset()
	var injectedFunc1 func(x int, y int) int
	var injectedFunc2 func(x float32) float32

	Require(&injectedFunc1)
	Require(&injectedFunc2)
	Provide(&injectedFunc1, f)
	Provide(&injectedFunc1, f3)

	errs := ResolveAll()
	if len(errs) != 2 {
		t.Fatal(errs)
	}
	assert.Nil(t, injectedFunc1)
	assert.Nil(t, injectedFunc2)

	fmt.Println(errs)

	if _, ok := errs[0].(*EMultipleFuncImplementations); !ok {
		t.Fatal(errs)
	}

	if _, ok := errs[1].(*EImplementationNotProvided); !ok {
		t.Fatal(errs)
	}
}

func TestErrorOnNonAssignableProvision(t *testing.T) {
	Reset()
	// target var non-ptr, required non-ptr
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	_, file, line, _ := runtime.Caller(0)
	Provide(injectedFunc, f)
	errs := ResolveAll()

	if e, ok := errs[0].(*EProvisionForNonAssignable); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, file, e.provisionPlace.file)
		assert.Equal(t, line+1, e.provisionPlace.line)
	} else {
		t.Fatal(errs)
	}
	assert.Nil(t, injectedFunc)

	Reset()
}

func TestErrorOnNonAssignableRequirement(t *testing.T) {
	Reset()
	
	var injectedFunc func(x int, y int) int

	_, file, line, _ := runtime.Caller(0)
	Require(injectedFunc)
	Provide(&injectedFunc, f)
	errs := ResolveAll()

	if e, ok := errs[0].(*ENonAssignableRequirement); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, file, e.req.file)
		assert.Equal(t, line+1, e.req.line)
	} else {
		t.Fatal(errs)
	}
	assert.Nil(t, injectedFunc)
}

func TestMatchReqAndImplByPointer(t *testing.T) {
	Reset()
	var injected1 func(x int, y int) int
	var injected2 func(x int, y int) int
	Require(&injected1)
	Require(&injected2)

	Provide(&injected1, f3)
	Provide(&injected2, f)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}

	assert.Equal(t, 6, injected1(2, 3))
	assert.Equal(t, 5, injected2(2, 3))
}

func TestErrorOnIncompatibleTypes(t *testing.T) {
	Reset()
	var injected func(x int, y int) int

	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&injected)
	_, implFile, implLine, _ := runtime.Caller(0)
	Provide(&injected, f2)

	errs := ResolveAll()

	if e, ok := errs[0].(*EIncompatibleTypesFunc); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, implFile, e.prov.file)
		assert.Equal(t, implLine+1, e.prov.line)
	} else {
		t.Fatal()
	}
	assert.Nil(t, injected)
}

func TestPackageUsedIfSomethigRequired(t *testing.T) {
	Reset()

	var injected1 func(x int, y int) int
	var injected2 func(x float32) float32

	Provide(&injected1, f)
	Provide(&injected2, f2)
	Require(&injected1)

	errs := ResolveAll()
	if errs != nil {
		t.Fatal()
	}
	assert.NotNil(t, injected1)
	assert.NotNil(t, injected2)
}

func TestPackageNotUsedIfNothingRequired(t *testing.T) {
	Reset()

	var injected1 func(x int, y int) int
	var injected2 func(x float32) float32

	Provide(&injected1, f)
	Provide(&injected2, f2)

	errs := ResolveAll()

	pc, _, _, _ := runtime.Caller(0)
	nameFull := runtime.FuncForPC(pc).Name()
	pkgName := nameFull[:strings.LastIndex(nameFull, ".")]
	if e, ok := errs[0].(*EPackageNotUsed); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, pkgName, e.pkgName)
	} else {
		t.Fatal()
	}
	assert.Nil(t, injected1)
	assert.Nil(t, injected2)
}

func TestDataInject(t *testing.T) {
	Reset()
	var injected map[string]int
	Provide(&injected, make(map[string]int))

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.NotNil(t, injected)
}

func TestProvideExtensionMapBasic(t *testing.T) {
	Reset()
	assert := assert.New(t)
	type bucketDef struct {
		Value string
	}

	type key struct {
		keyValue int
	}

	initedMap := make(map[string]int)
	var bucketDefsPtr map[string]*bucketDef
	var bucketDefs map[string]bucketDef
	var bucketDefsByKey map[key]bucketDef

	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	Provide(&bucketDefs, map[string]bucketDef{})
	Provide(&bucketDefsByKey, map[key]bucketDef{})

	var bucketServicePtr = &bucketDef{Value: "val"}
	var bucketService = bucketDef{Value: "val"}
	ProvideKeyValue(&initedMap, "str", 1)
	ProvideKeyValue(&bucketDefs, "key", bucketService)
	ProvideKeyValue(&bucketDefsPtr, "key", bucketServicePtr)
	ProvideKeyValue(&bucketDefsByKey, key{42}, bucketService)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(1, len(initedMap))
	assert.Equal(1, initedMap["str"])
	assert.Equal(1, len(bucketDefs))
	assert.Equal(bucketService, bucketDefs["key"])
	assert.Equal(1, len(bucketDefsPtr))
	assert.Equal(bucketServicePtr, bucketDefsPtr["key"])
	assert.Equal(1, len(bucketDefsByKey))
	assert.Equal(bucketService, bucketDefsByKey[key{42}])

	Reset()
	assert.Nil(bucketDefsPtr)
	assert.Nil(bucketDefs)
	assert.Nil(bucketDefsByKey)
	assert.NotNil(initedMap)
}

func TestProvideExtensionMapErrorOnNoProvideForNil(t *testing.T) {
	Reset()
	var myMap map[string]int

	ProvideKeyValue(&myMap, "str", 1)

	errs := ResolveAll()
	if _, ok := errs[0].(*EImplementationNotProvided); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal()
	}
}

func TestProvideExtensionMapErrorOnProvideForNonNil(t *testing.T) {
	Reset()
	myMap := make(map[string]int)

	Provide(&myMap, map[string]int{})
	ProvideKeyValue(&myMap, "str", 1)

	errs := ResolveAll()
	if _, ok := errs[0].(*EImplementationProvidedForNonNil); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal()
	}
}

func TestProvideExtensionMapErrorOnIncompatibleTypesKey(t *testing.T) {
	Reset()
	assert := assert.New(t)
	type bucketDef struct {
		Key string
	}
	type anotherType struct {
		Value string
	}

	var bucketDefsPtr map[string]*bucketDef

	var bucketServicePtr = &bucketDef{Key: "val"}
	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	_, provFile, provLine, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefsPtr, 1, bucketServicePtr)

	errs := ResolveAll()
	assert.Nil(bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageKey); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		assert.Equal(provFile, e.prov.file)
		assert.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnIncompatibleTypesValue(t *testing.T) {
	Reset()
	assert := assert.New(t)
	type bucketDef struct {
		Key string
	}
	type anotherType struct {
		Value string
	}

	var bucketDefsPtr map[string]*bucketDef

	var bucketServicePtr = &anotherType{Value: "val"}
	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	_, provFile, provLine, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefsPtr, "key", bucketServicePtr)

	errs := ResolveAll()
	assert.Nil(bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageValue); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		assert.Equal(provFile, e.prov.file)
		assert.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnIncompatibleTypesProvide(t *testing.T) {
	Reset()
	assert := assert.New(t)
	type bucketDef struct {
		Key string
	}

	var bucketDefsPtr map[string]bucketDef

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&bucketDefsPtr, map[string]*bucketDef{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EIncompatibleTypesStorageImpl); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		assert.Equal(provFile, e.prov.file)
		assert.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnMultipleProvisions(t *testing.T) {
	Reset()
	assert := assert.New(t)
	var m map[string]int

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&m, map[string]int{})
	Provide(&m, map[string]int{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EMultipleStorageImplementations); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Len(e.provs, 2)
		assert.Equal(provFile, e.provs[0].file)
		assert.Equal(provLine+1, e.provs[0].line)
		assert.Equal(provFile, e.provs[1].file)
		assert.Equal(provLine+2, e.provs[1].line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnAppendSliceIncompatibleTypes(t *testing.T) {
	Reset()
	type bucketDef struct {
		Key string
	}
	type anotherType struct {
		Value string
	}

	var bucketDefsPtr map[string][]*bucketDef

	var bucketServicePtr = &anotherType{Value: "val"}
	Provide(&bucketDefsPtr, map[string][]*bucketDef{})
	_, provFile, provLine, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefsPtr, "key", bucketServicePtr)

	errs := ResolveAll()
	assert.Nil(t, bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageValue); ok && len(errs) == 1 {
		fmt.Println(e)
		assert.Equal(t, reflect.TypeOf(bucketDefsPtr), e.reqType)
		assert.Equal(t, provFile, e.prov.file)
		assert.Equal(t, provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnMultipleValuesPerKey(t *testing.T) {
	Reset()
	type bucketDef struct {
		Value string
	}

	var bucketDefs map[string]*bucketDef

	Provide(&bucketDefs, map[string]*bucketDef{})
	Require(&bucketDefs)
	_, implFileF, implLineF, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefs, "key", &bucketDef{"val1"})
	_, implFileF2, implLineF2, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefs, "key", &bucketDef{"val2"})

	errs := ResolveAll()
	assert.Nil(t, bucketDefs)

	if e, ok := errs[0].(*EMultipleValues); ok && len(errs) == 1 {
		fmt.Println(e)
		assert.Equal(t, 2, len(e.provs))
		assert.Equal(t, implLineF+1, e.provs[0].line)
		assert.Equal(t, implFileF, e.provs[0].file)
		assert.Equal(t, implLineF2+1, e.provs[1].line)
		assert.Equal(t, implFileF2, e.provs[1].file)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionSliceBasic(t *testing.T) {
	Reset()
	assert := assert.New(t)
	type bucketDef struct {
		Value string
	}
	initedSlice := []int{42}
	var mySlice []string
	var bucketDefs []bucketDef
	var bucketDefsPtr []*bucketDef

	Provide(&mySlice, make([]string, 0))
	Provide(&bucketDefs, make([]bucketDef, 0))
	Provide(&bucketDefsPtr, make([]*bucketDef, 0))

	var bucketServicePtr = &bucketDef{Value: "val"}
	var bucketService = bucketDef{Value: "val"}
	ProvideSliceElement(&initedSlice, 43)
	ProvideSliceElement(&initedSlice, []int{44})
	ProvideSliceElement(&mySlice, "str1")
	ProvideSliceElement(&mySlice, []string{"str2"})
	ProvideSliceElement(&bucketDefs, bucketService)
	ProvideSliceElement(&bucketDefs, []bucketDef{bucketService})
	ProvideSliceElement(&bucketDefsPtr, []*bucketDef{bucketServicePtr})
	ProvideSliceElement(&bucketDefsPtr, bucketServicePtr)

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}

	assert.Len(initedSlice, 3)
	assert.Len(mySlice, 2)
	assert.Len(bucketDefs, 2)
	assert.Len(bucketDefsPtr, 2)
	assert.Equal(42, initedSlice[0])
	assert.Equal(43, initedSlice[1])
	assert.Equal(44, initedSlice[2])
	assert.Equal("str1", mySlice[0])
	assert.Equal("str2", mySlice[1])
	assert.Equal("val", bucketDefs[0].Value)
	assert.Equal("val", bucketDefs[1].Value)
	assert.Equal("val", bucketDefsPtr[0].Value)
	assert.Equal("val", bucketDefsPtr[1].Value)

	Reset()
	assert.Nil(mySlice)
	assert.Nil(bucketDefs)
	assert.Nil(bucketDefsPtr)
	assert.NotNil(initedSlice)
}

func TestProvideExtensionSliceErrorOnNoProvided(t *testing.T) {
	Reset()
	var mySlice []string

	ProvideSliceElement(&mySlice, "str")

	errs := ResolveAll()
	if _, ok := errs[0].(*EImplementationNotProvided); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal()
	}
}

func TestProvideExtensionSliceErrorOnProvideForNonNil(t *testing.T) {
	Reset()
	mySlice := []string{}
	mySliceImpl := []string{}

	Provide(&mySlice, mySliceImpl)

	errs := ResolveAll()
	if _, ok := errs[0].(*EImplementationProvidedForNonNil); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal()
	}
}

func TestProvideExtensionSliceErrorOnIncompatibleTypesSliceElement(t *testing.T) {
	Reset()
	var mySlice []string
	Provide(&mySlice, make([]string, 0))
	_, provFile, provLine, _ := runtime.Caller(0)
	ProvideSliceElement(&mySlice, 1)
	errs := ResolveAll()
	if e, ok := errs[0].(*EIncompatibleTypesStorageValue); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, reflect.TypeOf(mySlice), e.reqType)
		assert.Equal(t, provLine+1, e.prov.line)
		assert.Equal(t, provFile, e.prov.file)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionSliceErrorOnIncompatibleTypesProvide(t *testing.T) {
	Reset()
	var mySlice []string

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&mySlice, make([]int, 0))

	errs := ResolveAll()
	if e, ok := errs[0].(*EIncompatibleTypesStorageImpl); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(t, reflect.TypeOf(mySlice), e.reqType)
		assert.Equal(t, provLine+1, e.prov.line)
		assert.Equal(t, provFile, e.prov.file)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionSliceErrorOnMultipleProvisions(t *testing.T) {
	Reset()
	assert := assert.New(t)
	var s []string

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&s, []string{})
	Provide(&s, []string{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EMultipleStorageImplementations); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Len(e.provs, 2)
		assert.Equal(provFile, e.provs[0].file)
		assert.Equal(provLine+1, e.provs[0].line)
		assert.Equal(provFile, e.provs[1].file)
		assert.Equal(provLine+2, e.provs[1].line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapAppendSlice(t *testing.T) {
	Reset()
	assert := assert.New(t)

	type bucketDef struct {
		Value string
	}

	type key struct {
		keyValue int
	}

	inited := map[string][]int{}
	var intMap map[string][]int
	var intMapPtr map[string][]*int
	var bucketDefsPtr map[string][]*bucketDef
	var bucketDefs map[string][]bucketDef
	var bucketDefsByKey map[key][]bucketDef

	var bucketServicePtr1 = &bucketDef{Value: "val1"}
	var bucketServicePtr2 = &bucketDef{Value: "val2"}
	var bucketService = bucketDef{Value: "val"}

	Provide(&intMap, map[string][]int{})
	Require(&intMap)
	Provide(&intMapPtr, map[string][]*int{})
	Require(&intMapPtr)
	Provide(&bucketDefsPtr, map[string][]*bucketDef{})
	Require(&bucketDefsPtr)
	Provide(&bucketDefs, map[string][]bucketDef{})
	Require(&bucketDefs)
	Provide(&bucketDefsByKey, map[key][]bucketDef{})
	Require(&bucketDefsByKey)
	inited["tmp"] = []int{42}

	val1 := 3
	val2 := 4
	ProvideKeyValue(&inited, "tmp", []int{44})
	ProvideKeyValue(&inited, "tmp", 43)
	ProvideKeyValue(&intMap, "str", 1)
	ProvideKeyValue(&intMap, "str", []int{2})
	ProvideKeyValue(&intMapPtr, "str", &val1)
	ProvideKeyValue(&intMapPtr, "str", []*int{&val2})
	ProvideKeyValue(&bucketDefs, "key", bucketService)
	ProvideKeyValue(&bucketDefs, "key", []bucketDef{bucketService})
	ProvideKeyValue(&bucketDefsPtr, "key", bucketServicePtr1)
	ProvideKeyValue(&bucketDefsPtr, "key", []*bucketDef{bucketServicePtr2})
	ProvideKeyValue(&bucketDefsByKey, key{42}, bucketService)
	ProvideKeyValue(&bucketDefsByKey, key{42}, []bucketDef{bucketService})

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Len(inited, 1)
	assert.Len(inited["tmp"], 3)
	assert.Len(intMap, 1)
	assert.Len(intMap["str"], 2)
	assert.Len(intMapPtr, 1)
	assert.Len(intMapPtr["str"], 2)
	assert.Len(bucketDefs, 1)
	assert.Len(bucketDefs["key"], 2)
	assert.Len(bucketDefsPtr, 1)
	assert.Len(bucketDefsPtr["key"], 2)
	assert.Len(bucketDefsByKey, 1)
	assert.Len(bucketDefsByKey[key{42}], 2)

	assert.Equal(42, inited["tmp"][0])
	assert.Equal(44, inited["tmp"][1])
	assert.Equal(43, inited["tmp"][2])
	assert.Equal(1, intMap["str"][0])
	assert.Equal(2, intMap["str"][1])
	assert.Equal(val1, *intMapPtr["str"][0])
	assert.Equal(val2, *intMapPtr["str"][1])
	assert.Equal(bucketService, bucketDefs["key"][0])
	assert.Equal(bucketService, bucketDefs["key"][1])
	assert.Equal(bucketServicePtr1, bucketDefsPtr["key"][0])
	assert.Equal(bucketServicePtr2, bucketDefsPtr["key"][1])
	assert.Equal(bucketService, bucketDefsByKey[key{42}][0])
	assert.Equal(bucketService, bucketDefsByKey[key{42}][1])
}

func TestMixedRequirementsTypes(t *testing.T) {
	Reset()
	assert := assert.New(t)
	var injectedFunc func(x int, y int) int
	Require(&injectedFunc)
	Provide(&injectedFunc, f)

	var bucketDefs map[string]int
	Provide(&bucketDefs, map[string]int{})
	Require(&bucketDefs)
	ProvideKeyValue(&bucketDefs, "key", 1)

	var mySlice []string
	Require(&mySlice)
	Provide(&mySlice, []string{})
	ProvideSliceElement(&mySlice, "str1")

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}

	assert.Equal(5, injectedFunc(3, 2))
	assert.Equal(1, len(bucketDefs))
	assert.Equal(1, bucketDefs["key"])
	assert.Len(mySlice, 1)
	assert.Equal("str1", mySlice[0])
}

func TestErrorOnResoveTwice(t *testing.T) {
	Reset()
	assert := assert.New(t)
	var injectedFunc func(x int, y int) int
	Require(&injectedFunc)
	Provide(&injectedFunc, f)

	_, resFile, resLine, _ := runtime.Caller(0)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}

	errs = ResolveAll()
	if e, ok := errs[0].(*EAlreadyResolved); ok && len(errs) == 1 {
		fmt.Println(errs)
		assert.Equal(resFile, e.resolvePlace.file)
		assert.Equal(resLine+1, e.resolvePlace.line)
	} else {
		t.Fatal()
	}
}

type TMyType uint16

func TestReturnCustomType(t *testing.T) {
	Reset()
	var injectedFunc func(ctx context.Context) TMyType
	Require(&injectedFunc)
	errs := ResolveAll()
	assert.True(t, len(errs) > 0)
	Reset()
	Require(&injectedFunc)
	Provide(&injectedFunc, f4)
	errs = ResolveAll()
	assert.True(t, len(errs) == 0, errs)
}

func f(x int, y int) int {
	return x + y
}

func f2(x float32) float32 {
	return x + 1
}

func f3(x int, y int) int {
	return x * y
}

func f4(ctx context.Context) TMyType {
	return TMyType(10)
}
