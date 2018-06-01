package godif

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Func1Type = func(x int, y int) int
type Func2Type = func(x float32) float32

//var Put func(ctx context.Context, key interface{}, value interface{})

// output to log, not to console

func TestBasic(t *testing.T) {
	Reset()
	var injectedFunc Func1Type
	//var tmp func(x int, y int) int

	// tmp() // show
	// injectedFunc() // does not show


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

func TestErrorOnMultipleImplementations(t *testing.T) {
	Reset()
	var injectedFunc Func1Type

	Require(&injectedFunc)
	ProvideByImpl(f)
	ProvideByImpl(f3)
	errs := ResolveAll()
	if errs == nil {
		t.Fatal()
	}
	switch errs[0].(type) {
	case *EMultipleImplementations:
		fmt.Println(errs[0])
	default:
		t.Fatal()
	}

	assert.Nil(t, injectedFunc)
}

func TestErrorOnNonAssignableRequirement(t *testing.T) {
	Reset()
	var injectedFunc *Func1Type

	Require(injectedFunc)
	Provide(injectedFunc, f)
	errs := ResolveAll()
	if errs == nil {
		t.Fatal()
	}
	switch errs[0].(type) {
	case *ENonAssignableRequirement:
		fmt.Println(errs[0])
	default:
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
