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
	Provide(&injectedFunc1, f2)

	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EMultipleImplementations); ok {
		fmt.Println(e)
		assert.Equal(t, 2, len(e.impls))
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, implLineF+1, e.impls[0].line)
		assert.Equal(t, implFileF, e.impls[0].file)
		assert.Equal(t, implLineF2+1, e.impls[1].line)
		assert.Equal(t, implFileF2, e.impls[1].file)
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
	Provide(&injectedFunc1, f2)

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
		t.Fatal()
	}
}

func TestErrorOnNonAssignableRequirementWrongKind(t *testing.T) {
	Reset()
	var injected *func(x int, y int) int

	_, file, line, _ := runtime.Caller(0)
	Require(f)
	Provide(&injected, f)
	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*ENonAssignableRequirement); ok {
		fmt.Println(e)
		assert.Equal(t, file, e.req.file)
		assert.Equal(t, line+1, e.req.line)
	} else {
		t.Fatal()
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
	assert.NotNil(t, errs)

	if e, ok := errs[0].(*EIncompatibleTypes); ok {
		fmt.Println(e)
		assert.Equal(t, reqFile, e.req.file)
		assert.Equal(t, reqLine+1, e.req.line)
		assert.Equal(t, implFile, e.impl.file)
		assert.Equal(t, implLine+1, e.impl.line)
	} else {
		t.Fatal()
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
