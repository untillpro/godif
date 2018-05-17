package godif

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Func1Type func(x int, y int) int

func TestSimple(t *testing.T) {
	var injectedFunc Func1Type
	RegisterImpl(f)
	Require(&injectedFunc)
	err := ResolveAll()
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, 5, f(3, 2))
}

func f(x int, y int) int {
	return x + y
}
