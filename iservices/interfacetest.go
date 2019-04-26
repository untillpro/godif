/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untillpro/godif"
)

var declareImplementation func()
var lastCtx context.Context

// TestImpl s.e.
func TestImpl(t *testing.T, declare func()) {
	declareImplementation = declare
	t.Run("Test_BasicUsage", testBasicUsage)
	t.Run("Test_FailedStart", testFailedStart)
}

func testBasicUsage(t *testing.T) {
	godif.Require(&Start)
	godif.Require(&Stop)

	// Declare passed implementation

	declareImplementation()

	// Provide own services

	s1 := &MyService{Name: "Service1"}
	s2 := &MyService{Name: "Service2"}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)

	// Resolve all

	errs := godif.ResolveAll()
	defer godif.Reset()
	require.Nil(t, errs, errs)

	// Start services

	var err error
	ctx := context.Background()
	ctx, err = Start(ctx)
	defer Stop(ctx)
	require.Nil(t, err)

	// Check service state

	assert.Equal(t, 1, s1.State)
	assert.Equal(t, 1, s2.State)

	//Make sure that value provided by service exist in ctx

	assert.True(t, lastCtx.Value(ctxKeyType("Service1")).(bool))
	assert.True(t, lastCtx.Value(ctxKeyType("Service2")).(bool))
	assert.Nil(t, lastCtx.Value(ctxKeyType("Service3")))

	// Stop services
	Stop(ctx)
	assert.Equal(t, 0, s1.State)
	assert.Equal(t, 0, s2.State)
}

func testFailedStart(t *testing.T) {
	godif.Require(&Start)
	godif.Require(&Stop)

	// Declare passed implementation

	declareImplementation()

	// Provide own services

	s1 := &MyService{Name: "Service1"}
	s2 := &MyService{Name: "Service2", Failstart: true}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)

	// Resolve all

	errs := godif.ResolveAll()
	defer godif.Reset()
	assert.Nil(t, errs)
	fmt.Println("errs=", errs)

	// Start services

	var err error
	fmt.Println("### Before Start")
	ctx := context.Background()
	ctx, err = Start(ctx)
	defer Stop(ctx)
	fmt.Println("### After Start")
	assert.NotNil(t, err)
	fmt.Println("err=", err)
	assert.True(t, strings.Contains(err.Error(), "Service2"))
	assert.False(t, strings.Contains(err.Error(), "Service1"))
	assert.Equal(t, 1, s1.State)
	assert.Equal(t, 0, s2.State)
}

// MyService for testing purposes
type MyService struct {
	Name      string
	State     int // 0 (stopped), 1 (started)
	Failstart bool
	CtxValue  interface{}
	Wg        *sync.WaitGroup
}

type ctxKeyType string

// Start s.e.
func (s *MyService) Start(ctx context.Context) (context.Context, error) {
	if s.Failstart {
		fmt.Println(s.Name, "Start fails")
		return ctx, errors.New(s.Name + ":" + "Start fails")
	}
	s.State++
	fmt.Println(s.Name, "Started")
	ctx = context.WithValue(ctx, ctxKeyType(s.Name), true)
	if nil != s.Wg {
		s.Wg.Done()
	}
	lastCtx = ctx
	return ctx, nil
}

// Stop s.e.
func (s *MyService) Stop(ctx context.Context) {
	s.State--
	fmt.Println(s.Name, "Stopped")
}

func (s *MyService) String() string {
	return "I'm service " + s.Name
}
