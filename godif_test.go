package godif

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Func1Type = func(x int, y int) int

func TestSimple(t *testing.T) {
	var injectedFunc Func1Type

	err := ResolveAll()
	assert.Nil(t, err)
	
	Require(&injectedFunc)
	err = ResolveAll()
	assert.NotNil(t, err)

	RegisterImpl(f)
	err = ResolveAll()
	assert.Nil(t, err)
	assert.Equal(t, 5, injectedFunc(3, 2))
}

func f(x int, y int) int {
	return x + y
}
