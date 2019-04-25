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
	"testing"
	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/untillpro/godif"
)

var declareImplementation func()
var lastCtx context.Context

// TestAll s.e.
func TestAll(t *testing.T, declare func()) {
	declareImplementation = declare
	t.Run("TestRun", testRun)
	t.Run("TestFailedStart", testFailedStart)
}

func testRun(t *testing.T) {

	godif.Require(&Start)
	godif.Require(&Stop)

	// Declare passed implementation

	declareImplementation()

	// Provide own services

	var wg sync.WaitGroup
	wg.Add(2)
	s1 := &myService{Name: "Service1", wg: &wg}
	s2 := &myService{Name: "Service2", wg: &wg}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)

	// Init and start services

	var err error
	done :=  make(chan string)
	go func(){
		fmt.Println("### Before Run")
		err = Run()
		fmt.Println("### After Run")		
		assert.Nil(t, err)
		assert.Equal(t, 0, s1.State)
		assert.Equal(t, 0, s2.State)
		done<-""
	}()

	defer func(){
		fmt.Println("### Before Terminate()")
		Terminate()
		fmt.Println("### After Terminate()")		
		<-done
		fmt.Println("### Done")
	}()

	fmt.Println("### Before wg.Wait()")
	wg.Wait()
	fmt.Println("### After wg.Wait()")

	fmt.Println("### Start testing ctx=", lastCtx)

	//2 means "started"
	assert.Equal(t, 1, s1.State)
	assert.Equal(t, 1, s2.State)

	//Make sure that value provided by service exist in ctx
	assert.True(t, lastCtx.Value(ctxKeyType("Service1")).(bool))
	assert.True(t, lastCtx.Value(ctxKeyType("Service2")).(bool))
	assert.Nil(t, lastCtx.Value(ctxKeyType("Service3")))
}

func testFailedStart(t *testing.T) {
	godif.Require(&Start)
	godif.Require(&Stop)

	// Declare passed implementation

	declareImplementation()

	// Provide own services

	s1 := &myService{Name: "Service1"}
	s2 := &myService{Name: "Service2", Failstart: true}
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

type myService struct {
	Name      string
	State     int // 0 (stopped), 1 (started)
	Failstart bool
	ctxValue  interface{}
	wg *sync.WaitGroup
}

type ctxKeyType string

func (s *myService) Start(ctx context.Context) (context.Context, error) {
	if s.Failstart {
		fmt.Println(s.Name, "Start fails")
		return ctx, errors.New(s.Name + ":" + "Start fails")
	}
	s.State++
	fmt.Println(s.Name, "Started")
	ctx = context.WithValue(ctx, ctxKeyType(s.Name), true)
	if nil != s.wg {
		s.wg.Done()
	}
	lastCtx = ctx
	return ctx, nil
}

func (s *myService) Stop(ctx context.Context) {
	s.State--
	fmt.Println(s.Name, "Stopped")
}

func (s *myService) String() string {
	return "I'm service " + s.Name
}