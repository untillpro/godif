/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicUsage(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Nil(t, injectedFunc)

	Require(&injectedFunc)
	errs = ResolveAll()
	if errs == nil {
		t.Fatal()
	}
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	fmt.Println(errs)

	switch errs[0].(type) {
	case *EImplementationNotProvided:
		fmt.Println(errs[0])
	default:
		t.Fatal()
	}
	assert.Nil(t, injectedFunc)

	Provide(&injectedFunc, f)
	errs = ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, injectedFunc(3, 2))

	Reset()
	assert.Nil(t, injectedFunc)
}

func TestExplicitTypeInject(t *testing.T) {
	type Func1Type = func(x int, y int) int
	Reset()
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
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EMultipleImplementations); ok {
		fmt.Println(e)
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

	fmt.Println(errs)

	switch errs[0].(type) {
	case *EMultipleImplementations:
	default:
		t.Fatal()
	}

	switch errs[1].(type) {
	case *EImplementationNotProvided:
	default:
		t.Fatal()
	}
}

func TestErrorOnNonAssignableRequirementNonPointer(t *testing.T) {
	Reset()
	var injectedFunc *func(x int, y int) int

	_, file, line, _ := runtime.Caller(0)
	Require(injectedFunc)
	Provide(injectedFunc, f)
	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*ENonAssignableRequirement); ok {
		fmt.Println(e)
		assert.Equal(t, file, e.req.file)
		assert.Equal(t, line+1, e.req.line)
	} else {
		t.Fatal(errs)
	}
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

func TestDoNotResolveAtAllOnAnyError(t *testing.T) {
	var injected1 func(x int, y int) int
	var injected2 func(x int, y int) int

	Require(&injected1)
	Require(&injected2)
	Provide(&injected1, f)

	errs := ResolveAll()
	assert.Equal(t, 1, len(errs))

	assert.Nil(t, injected1)
	assert.Nil(t, injected2)
}

func TestErrorOnIncompatibleTypes(t *testing.T) {
	Reset()
	var injected func(x int, y int) int

	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&injected)
	_, implFile, implLine, _ := runtime.Caller(0)
	Provide(&injected, f2)

	errs := ResolveAll()
	assert.Equal(t, 1, len(errs))

	if e, ok := errs[0].(*EIncompatibleTypes); ok {
		fmt.Println(e)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, implFile, e.prov.file)
		assert.Equal(t, implLine+1, e.prov.line)
	} else {
		t.Fatal()
	}
}

func TestErrorOnProvidedButNotUsed(t *testing.T) {
	Reset()
	var injected func(x int, y int) int

	_, implFile, implLine, _ := runtime.Caller(0)
	Provide(&injected, f)

	errs := ResolveAll()
	assert.Equal(t, 1, len(errs))

	fmt.Println(errs[0])

	if e, ok := errs[0].(*EProvidedNotUsed); ok {
		assert.Equal(t, implFile, e.prov.file)
		assert.Equal(t, implLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}

	Require(&injected)
	errs = ResolveAll()

	assert.Nil(t, errs)
}

func TestDataInject(t *testing.T) {
	Reset()
	var injected map[string]int
	Require(&injected)
	Provide(&injected, make(map[string]int))

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.NotNil(t, injected)
}

func TestErrorOnIncompatibleTypesDataInject(t *testing.T) {
	Reset()
	var injected map[string]int
	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&injected)
	_, provFile, provLine, _ := runtime.Caller(0)
	Provide(&injected, make([]int, 0))

	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EIncompatibleTypes); ok {
		fmt.Println(e)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, provFile, e.prov.file)
		assert.Equal(t, provLine+1, e.prov.line)
	} else {
		t.Fatal()
	}
}

func TestProvideMapValue(t *testing.T) {
	Reset()
	type bucketDef struct {
		Value string
	}

	type key struct {
		keyValue int
	}

	var bucketDefsPtr map[string]*bucketDef
	var bucketDefs map[string]bucketDef
	var bucketDefsByKey map[key]bucketDef

	var bucketServicePtr = &bucketDef{Value: "val"}
	var bucketService = bucketDef{Value: "val"}
	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	Require(&bucketDefsPtr)
	Provide(&bucketDefs, map[string]bucketDef{})
	Require(&bucketDefs)
	Provide(&bucketDefsByKey, map[key]bucketDef{})
	Require(&bucketDefsByKey)

	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Empty(t, bucketDefs)
	assert.Empty(t, bucketDefsPtr)
	assert.Empty(t, bucketDefsByKey)

	ProvideMapValue(&bucketDefs, "key", bucketService)
	ProvideMapValue(&bucketDefsPtr, "key", bucketServicePtr)
	ProvideMapValue(&bucketDefsByKey, key{42}, bucketService)
	errs = ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 1, len(bucketDefs))
	assert.Equal(t, bucketService, bucketDefs["key"])
	assert.Equal(t, 1, len(bucketDefsPtr))
	assert.Equal(t, bucketServicePtr, bucketDefsPtr["key"])
	assert.Equal(t, 1, len(bucketDefsByKey))
	assert.Equal(t, bucketService, bucketDefsByKey[key{42}])
}

func TestProvideMapValueIncompatibleTypes(t *testing.T) {
	Reset()
	type bucketDef struct {
		Key string
	}
	type anotherType struct {
		Value string
	}

	var bucketDefsPtr map[string]*bucketDef

	var bucketServicePtr = &anotherType{Value: "val"}
	Provide(&bucketDefsPtr, map[string]*bucketDef{})
	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&bucketDefsPtr)
	_, provFile, provLine, _ := runtime.Caller(0)
	ProvideMapValue(&bucketDefsPtr, "key", bucketServicePtr)

	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EIncompatibleTypes); ok {
		fmt.Println(e)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, provFile, e.prov.file)
		assert.Equal(t, provLine+1, e.prov.line)
	} else {
		t.Fatal(errs)
	}
}

func TestProvideMapValueErrorOnMultipleValuesPerKey(t *testing.T) {
	Reset()
	type bucketDef struct {
		Value string
	}

	var bucketDefs map[string]*bucketDef

	Provide(&bucketDefs, map[string]*bucketDef{})
	_, reqFile, reqLine, _ := runtime.Caller(0)
	Require(&bucketDefs)
	_, implFileF, implLineF, _ := runtime.Caller(0)
	ProvideMapValue(&bucketDefs, "key", &bucketDef{"val1"})
	_, implFileF2, implLineF2, _ := runtime.Caller(0)
	ProvideMapValue(&bucketDefs, "key", &bucketDef{"val2"})

	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EMultipleValues); ok {
		fmt.Println(e)
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
