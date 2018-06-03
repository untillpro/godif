package godif

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Func1Type = func(x int, y int) int
type Func2Type = func(x float32) float32

func TestExplicitType(t *testing.T) {
	Reset()
	var injectedFunc Func1Type

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

	ProvideByImpl(f)
	errs = ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, injectedFunc(3, 2))

	Reset()
	assert.Nil(t, injectedFunc)
}

func TestImplicitType(t *testing.T) {
	var inject func(x int, y int) int
	Require(&inject)
	ProvideByImpl(f)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, inject(2, 3))
}

func TestMultipleImplementationsError(t *testing.T) {
	Reset()
	var injectedFunc1 Func1Type

	Require(&injectedFunc1)
	_, fileF, lineF, _ := runtime.Caller(0); 
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
		assert.Equal(t, lineF + 1, e.impls[0].line)
		assert.Equal(t, fileF, e.impls[0].file)
		assert.Equal(t, lineF2 + 1, e.impls[1].line)
		assert.Equal(t, fileF2, e.impls[1].file)
	} else {
		t.Fatal(errs)
	}
}

func TestMultipleErrorsOnResolve(t *testing.T) {
	Reset()
	var injectedFunc1 Func1Type
	var injectedFunc2 Func2Type

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
		fmt.Println(errs[0])
	default:
		t.Fatal()
	}

	switch errs[1].(type) {
	case *EImplementationNotProvided:
		fmt.Println(errs[1])
	default:
		t.Fatal()
	}

}

func TestProvideByVar(t *testing.T) {
	Reset()
	var injectedFunc1 Func1Type
	var injectedFunc2 Func2Type

	Require(&injectedFunc1)
	Require(&injectedFunc2)
	Provide(&injectedFunc1, f)
	Provide(&injectedFunc2, f2)
	errs := ResolveAll()
	if errs != nil {
		t.Fatal(errs)
	}
	assert.Equal(t, 5, injectedFunc1(3, 2))
	assert.Equal(t, float32(3.5), injectedFunc2(2.5))
}

func TestErrorOnNonAssignableRequirement(t *testing.T) {
	Reset()
	var injectedFunc *Func1Type

	_, file, line, _ := runtime.Caller(0); 
	Require(injectedFunc)
	Provide(injectedFunc, f)
	errs := ResolveAll()
	if len(errs) != 1 {
		t.Fatal(errs)
	}

	if e, ok := errs[0].(*ENonAssignableRequirement); ok {
		fmt.Println(e)
		assert.Equal(t, file, e.requirement.file)
		assert.Equal(t, line + 1, e.requirement.line)
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
