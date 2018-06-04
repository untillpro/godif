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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromScratchImplicitTypeInject(t *testing.T) {
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

	Require(&injectedFunc1)
	_, fileF, lineF, _ := runtime.Caller(0)
	Provide(&injectedFunc1, f)
	_, fileF2, lineF2, _ := runtime.Caller(0)
	Provide(&injectedFunc1, f2)

	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*EMultipleImplementations); ok {
		fmt.Println(e)
		assert.Equal(t, reflect.TypeOf(&injectedFunc1), e.Type)
		assert.Equal(t, 2, len(e.impls))
		assert.Equal(t, lineF+1, e.impls[0].line)
		assert.Equal(t, fileF, e.impls[0].file)
		assert.Equal(t, lineF2+1, e.impls[1].line)
		assert.Equal(t, fileF2, e.impls[1].file)
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
		assert.Equal(t, file, e.requirement.file)
		assert.Equal(t, line+1, e.requirement.line)
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
		assert.Equal(t, file, e.requirement.file)
		assert.Equal(t, line+1, e.requirement.line)
	} else {
		t.Fatal()
	}
}

func TestMatchByPointer(t *testing.T) {
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

func f(x int, y int) int {
	return x + y
}

func f2(x float32) float32 {
	return x + 1
}

func f3(x int, y int) int {
	return x * y
}
