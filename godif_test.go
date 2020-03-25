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

	"github.com/stretchr/testify/require"
)

func TestFuncBasic(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	Provide(&injectedFunc, f)

	errs := ResolveAll()
	require.Nil(t, errs)
	require.Equal(t, 5, injectedFunc(3, 2))

	Reset()
	require.Nil(t, injectedFunc)
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

	require.Nil(t, injectedFunc)
}

func TestExplicitTypeInject(t *testing.T) {
	Reset()
	type Func1Type = func(x int, y int) int
	var inject Func1Type

	Require(&inject)
	Provide(&inject, f)
	errs := ResolveAll()
	require.Nil(t, errs)
	require.Equal(t, 5, inject(2, 3))
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
		require.Equal(t, 2, len(e.provs))
		require.Equal(t, reqLine+1, e.req.line)
		require.Equal(t, reqFile, e.req.file)
		require.Equal(t, implLineF+1, e.provs[0].line)
		require.Equal(t, implFileF, e.provs[0].file)
		require.Equal(t, implLineF2+1, e.provs[1].line)
		require.Equal(t, implFileF2, e.provs[1].file)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, injectedFunc1)
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
	require.Nil(t, injectedFunc1)
	require.Nil(t, injectedFunc2)

	fmt.Println(errs)

	if _, ok := errs[0].(*EMultipleFuncImplementations); !ok {
		t.Fatal(errs)
	}

	if _, ok := errs[1].(*EImplementationNotProvided); !ok {
		t.Fatal(errs)
	}
}

func TestErrorOnNonAssignableVarOnProvideFunc(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	_, file, line, _ := runtime.Caller(0)
	Provide(injectedFunc, f)
	errs := ResolveAll()

	if e, ok := errs[0].(*EProvisionForNonAssignable); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.provisionPlace.file)
		require.Equal(t, line+1, e.provisionPlace.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, injectedFunc)

}
func TestErrorOnNonAssignableVarOnProvideMap(t *testing.T) {
	Reset()

	var bucketDefs map[string]string
	Require(&bucketDefs)
	_, file, line, _ := runtime.Caller(0)
	Provide(bucketDefs, map[string]string{})
	errs := ResolveAll()

	if e, ok := errs[0].(*EProvisionForNonAssignable); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.provisionPlace.file)
		require.Equal(t, line+1, e.provisionPlace.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, bucketDefs)
}

func TestErrorOnNonAssignableVarOnProvideKeyValue(t *testing.T) {
	Reset()

	var intMap map[string][]int

	_, file, line, _ := runtime.Caller(0)
	ProvideKeyValue(intMap, "str", 1)

	errs := ResolveAll()
	if e, ok := errs[0].(*EProvisionForNonAssignable); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.provisionPlace.file)
		require.Equal(t, line+1, e.provisionPlace.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, intMap)
}

func TestErrorOnNonAssignableVarOnProvideSliceElement(t *testing.T) {
	Reset()

	var mySlice []string
	Provide(&mySlice, make([]string, 0))

	_, file, line, _ := runtime.Caller(0)
	ProvideSliceElement(mySlice, "str1")
	errs := ResolveAll()
	if e, ok := errs[0].(*EProvisionForNonAssignable); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.provisionPlace.file)
		require.Equal(t, line+1, e.provisionPlace.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, mySlice)

}

func TestErrorOnNonAssignableRequirementFunc(t *testing.T) {
	Reset()

	var injectedFunc func(x int, y int) int

	_, file, line, _ := runtime.Caller(0)
	Require(injectedFunc)
	errs := ResolveAll()

	if e, ok := errs[0].(*ENonAssignableRequirement); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.req.file)
		require.Equal(t, line+1, e.req.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, injectedFunc)
}

func TestErrorOnNonAssignableRequirementNonFunc(t *testing.T) {
	Reset()

	var bucketDefs map[string]string

	_, file, line, _ := runtime.Caller(0)
	Require(bucketDefs)
	errs := ResolveAll()

	if e, ok := errs[0].(*ENonAssignableRequirement); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, file, e.req.file)
		require.Equal(t, line+1, e.req.line)
	} else {
		t.Fatal(errs)
	}
	require.Nil(t, bucketDefs)
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
	require.Nil(t, errs)

	require.Equal(t, 6, injected1(2, 3))
	require.Equal(t, 5, injected2(2, 3))
}

func TestErrorOnIncompatibleTypesFunc(t *testing.T) {
	Reset()
	var injected func(x int, y int) int

	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&injected)
	_, implFile, implLine, _ := runtime.Caller(0)
	Provide(&injected, f2)

	errs := ResolveAll()

	if e, ok := errs[0].(*EIncompatibleTypesFunc); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(t, reqFile, e.req.file)
		require.Equal(t, reqLine+1, e.req.line)
		require.Equal(t, implFile, e.prov.file)
		require.Equal(t, implLine+1, e.prov.line)
	} else {
		t.Fatal()
	}
	require.Nil(t, injected)
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
	require.NotNil(t, injected1)
	require.Nil(t, injected2) // not required
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
		require.Equal(t, pkgName, e.pkgName)
	} else {
		t.Fatal()
	}
	require.Nil(t, injected1)
	require.Nil(t, injected2)
}

func TestDataInject(t *testing.T) {
	Reset()
	var injected map[string]int
	Provide(&injected, make(map[string]int))
	Require(&injected)

	errs := ResolveAll()
	require.Nil(t, errs)
	require.NotNil(t, injected)
}

func TestProvideExtensionMapBasic(t *testing.T) {
	Reset()
	require := require.New(t)
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

	// ProvideKeyValue() and Provide() called -> consider implicitly required
	ProvideKeyValue(&initedMap, "str", 1)
	ProvideKeyValue(&bucketDefs, "key", bucketService)
	ProvideKeyValue(&bucketDefsPtr, "key", bucketServicePtr)
	ProvideKeyValue(&bucketDefsByKey, key{42}, bucketService)
	errs := ResolveAll()
	require.Nil(errs)
	require.Equal(1, len(initedMap))
	require.Equal(1, initedMap["str"])
	require.Equal(1, len(bucketDefs))
	require.Equal(bucketService, bucketDefs["key"])
	require.Equal(1, len(bucketDefsPtr))
	require.Equal(bucketServicePtr, bucketDefsPtr["key"])
	require.Equal(1, len(bucketDefsByKey))
	require.Equal(bucketService, bucketDefsByKey[key{42}])

	Reset()
	require.Nil(bucketDefsPtr)
	require.Nil(bucketDefs)
	require.Nil(bucketDefsByKey)
	require.NotNil(initedMap)

	// test provided but no data provided -> not required -> no implementation
	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	errs = ResolveAll()
	require.Nil(errs)
	require.Nil(bucketDefsPtr)
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
	require := require.New(t)
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
	require.Nil(bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageKey); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		require.Equal(provFile, e.prov.file)
		require.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnIncompatibleTypesValue(t *testing.T) {
	Reset()
	require := require.New(t)
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
	require.Nil(bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageValue); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		require.Equal(provFile, e.prov.file)
		require.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnIncompatibleTypesProvide(t *testing.T) {
	Reset()
	require := require.New(t)
	type bucketDef struct {
		Key string
	}

	var bucketDefsPtr map[string]bucketDef

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&bucketDefsPtr, map[string]*bucketDef{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EIncompatibleTypesStorageImpl); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(reflect.TypeOf(bucketDefsPtr), e.reqType)
		require.Equal(provFile, e.prov.file)
		require.Equal(provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapErrorOnMultipleProvisions(t *testing.T) {
	Reset()
	require := require.New(t)
	var m map[string]int

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&m, map[string]int{})
	Provide(&m, map[string]int{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EMultipleStorageImplementations); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Len(e.provs, 2)
		require.Equal(provFile, e.provs[0].file)
		require.Equal(provLine+1, e.provs[0].line)
		require.Equal(provFile, e.provs[1].file)
		require.Equal(provLine+2, e.provs[1].line)
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
	require.Nil(t, bucketDefsPtr)

	if e, ok := errs[0].(*EIncompatibleTypesStorageValue); ok && len(errs) == 1 {
		fmt.Println(e)
		require.Equal(t, reflect.TypeOf(bucketDefsPtr), e.reqType)
		require.Equal(t, provFile, e.prov.file)
		require.Equal(t, provLine+1, e.prov.line)
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
	_, implFileF, implLineF, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefs, "key", &bucketDef{"val1"})
	_, implFileF2, implLineF2, _ := runtime.Caller(0)
	ProvideKeyValue(&bucketDefs, "key", &bucketDef{"val2"})

	errs := ResolveAll()
	require.Nil(t, bucketDefs)

	if e, ok := errs[0].(*EMultipleValues); ok && len(errs) == 1 {
		fmt.Println(e)
		require.Equal(t, 2, len(e.provs))
		require.Equal(t, implLineF+1, e.provs[0].line)
		require.Equal(t, implFileF, e.provs[0].file)
		require.Equal(t, implLineF2+1, e.provs[1].line)
		require.Equal(t, implFileF2, e.provs[1].file)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionSliceBasic(t *testing.T) {
	Reset()
	require := require.New(t)
	type bucketDef struct {
		Value string
	}
	initedSlice := []int{42}
	var mySlice []string
	var bucketDefs []bucketDef
	var bucketDefsPtr []*bucketDef
	var bucketServicePtr = &bucketDef{Value: "val"}
	var bucketService = bucketDef{Value: "val"}

	// implementation is not equired for nil slices
	ProvideSliceElement(&initedSlice, 43)
	ProvideSliceElement(&initedSlice, []int{44})
	ProvideSliceElement(&mySlice, "str1")
	ProvideSliceElement(&mySlice, []string{"str2"})
	ProvideSliceElement(&bucketDefs, bucketService)
	ProvideSliceElement(&bucketDefs, []bucketDef{bucketService})
	ProvideSliceElement(&bucketDefsPtr, []*bucketDef{bucketServicePtr})
	ProvideSliceElement(&bucketDefsPtr, bucketServicePtr)

	errs := ResolveAll()
	require.Nil(errs)

	require.Len(initedSlice, 3)
	require.Len(mySlice, 2)
	require.Len(bucketDefs, 2)
	require.Len(bucketDefsPtr, 2)
	require.Equal(42, initedSlice[0])
	require.Equal(43, initedSlice[1])
	require.Equal(44, initedSlice[2])
	require.Equal("str1", mySlice[0])
	require.Equal("str2", mySlice[1])
	require.Equal("val", bucketDefs[0].Value)
	require.Equal("val", bucketDefs[1].Value)
	require.Equal("val", bucketDefsPtr[0].Value)
	require.Equal("val", bucketDefsPtr[1].Value)

	// non-provided slices are not nilled on reset
	Reset()
	require.NotNil(mySlice)
	require.NotNil(bucketDefs)
	require.NotNil(bucketDefsPtr)
	require.NotNil(initedSlice)

	// provided not required slices are nilled on reset
	mySlice = nil
	bucketDefs = nil
	bucketDefsPtr = nil
	Provide(&mySlice, make([]string, 0))
	Provide(&bucketDefs, make([]bucketDef, 0))
	Provide(&bucketDefsPtr, make([]*bucketDef, 0))
	ProvideSliceElement(&initedSlice, 43)
	ProvideSliceElement(&mySlice, "str1")
	ProvideSliceElement(&bucketDefs, bucketService)
	ProvideSliceElement(&bucketDefsPtr, bucketServicePtr)
	errs = ResolveAll()
	require.Nil(errs)
	Reset()
	require.Nil(mySlice)
	require.Nil(bucketDefs)
	require.Nil(bucketDefsPtr)
	require.NotNil(initedSlice)

	// provided required slices are nilled on reset
	mySlice = nil
	bucketDefs = nil
	bucketDefsPtr = nil
	Provide(&mySlice, make([]string, 0))
	Provide(&bucketDefs, make([]bucketDef, 0))
	Provide(&bucketDefsPtr, make([]*bucketDef, 0))
	Require(&mySlice)
	Require(&bucketDefs)
	Require(&bucketDefsPtr)
	ProvideSliceElement(&initedSlice, 43)
	ProvideSliceElement(&mySlice, "str1")
	ProvideSliceElement(&bucketDefs, bucketService)
	ProvideSliceElement(&bucketDefsPtr, bucketServicePtr)
	errs = ResolveAll()
	require.Nil(errs)
	Reset()
	require.Nil(mySlice)
	require.Nil(bucketDefs)
	require.Nil(bucketDefsPtr)
	require.NotNil(initedSlice)

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
		require.Equal(t, reflect.TypeOf(mySlice), e.reqType)
		require.Equal(t, provLine+1, e.prov.line)
		require.Equal(t, provFile, e.prov.file)
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
		require.Equal(t, reflect.TypeOf(mySlice), e.reqType)
		require.Equal(t, provLine+1, e.prov.line)
		require.Equal(t, provFile, e.prov.file)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionSliceErrorOnMultipleProvisions(t *testing.T) {
	Reset()
	require := require.New(t)
	var s []string

	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&s, []string{})
	Provide(&s, []string{})

	errs := ResolveAll()

	if e, ok := errs[0].(*EMultipleStorageImplementations); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Len(e.provs, 2)
		require.Equal(provFile, e.provs[0].file)
		require.Equal(provLine+1, e.provs[0].line)
		require.Equal(provFile, e.provs[1].file)
		require.Equal(provLine+2, e.provs[1].line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideExtensionMapAppendSlice(t *testing.T) {
	Reset()
	require := require.New(t)

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
	Provide(&intMapPtr, map[string][]*int{})
	Provide(&bucketDefsPtr, map[string][]*bucketDef{})
	Provide(&bucketDefs, map[string][]bucketDef{})
	Provide(&bucketDefsByKey, map[key][]bucketDef{})
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
	require.Nil(errs)
	require.Len(inited, 1)
	require.Len(inited["tmp"], 3)
	require.Len(intMap, 1)
	require.Len(intMap["str"], 2)
	require.Len(intMapPtr, 1)
	require.Len(intMapPtr["str"], 2)
	require.Len(bucketDefs, 1)
	require.Len(bucketDefs["key"], 2)
	require.Len(bucketDefsPtr, 1)
	require.Len(bucketDefsPtr["key"], 2)
	require.Len(bucketDefsByKey, 1)
	require.Len(bucketDefsByKey[key{42}], 2)

	require.Equal(42, inited["tmp"][0])
	require.Equal(44, inited["tmp"][1])
	require.Equal(43, inited["tmp"][2])
	require.Equal(1, intMap["str"][0])
	require.Equal(2, intMap["str"][1])
	require.Equal(val1, *intMapPtr["str"][0])
	require.Equal(val2, *intMapPtr["str"][1])
	require.Equal(bucketService, bucketDefs["key"][0])
	require.Equal(bucketService, bucketDefs["key"][1])
	require.Equal(bucketServicePtr1, bucketDefsPtr["key"][0])
	require.Equal(bucketServicePtr2, bucketDefsPtr["key"][1])
	require.Equal(bucketService, bucketDefsByKey[key{42}][0])
	require.Equal(bucketService, bucketDefsByKey[key{42}][1])
}

func TestMixedRequirementsTypes(t *testing.T) {
	Reset()
	require := require.New(t)
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
	require.Nil(errs)

	require.Equal(5, injectedFunc(3, 2))
	require.Equal(1, len(bucketDefs))
	require.Equal(1, bucketDefs["key"])
	require.Len(mySlice, 1)
	require.Equal("str1", mySlice[0])
}

func TestErrorOnResoveTwice(t *testing.T) {
	Reset()
	require := require.New(t)
	var injectedFunc func(x int, y int) int
	Require(&injectedFunc)
	Provide(&injectedFunc, f)

	_, resFile, resLine, _ := runtime.Caller(0)
	errs := ResolveAll()
	require.Nil(errs)

	errs = ResolveAll()
	if e, ok := errs[0].(*EAlreadyResolved); ok && len(errs) == 1 {
		fmt.Println(errs)
		require.Equal(resFile, e.resolvePlace.file)
		require.Equal(resLine+1, e.resolvePlace.line)
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
	require.True(t, len(errs) > 0)
	Reset()
	Require(&injectedFunc)
	Provide(&injectedFunc, f4)
	errs = ResolveAll()
	require.True(t, len(errs) == 0, errs)
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
