/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/untillpro/godif"
	"github.com/untillpro/godif/iservices"
)

func Test_BasicUsage(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(2)

	// Declare two services

	s1 := &iservices.MyService{Name: "Service1", Wg: &wg}
	s2 := &iservices.MyService{Name: "Service2", Wg: &wg}
	godif.ProvideSliceElement(&iservices.Services, s1)
	godif.ProvideSliceElement(&iservices.Services, s2)

	// Terminate when all services started

	go func() {
		wg.Wait()
		Terminate()
	}()

	// Run waits for Terminate() or SIGTERM
	err := Run()
	require.Nil(t, err, err)

}

func Test_FailedStart(t *testing.T) {
	s1 := &iservices.MyService{Name: "Service1"}
	s2 := &iservices.MyService{Name: "Service2", Failstart: true}
	godif.ProvideSliceElement(&iservices.Services, s1)
	godif.ProvideSliceElement(&iservices.Services, s2)
	err := Run()
	require.NotNil(t, err, err)
	require.Equal(t, 0, s1.State)
	require.Equal(t, 0, s2.State)
}

var Missed func()

func Test_FailedResolve(t *testing.T) {
	godif.Require(&Missed)
	err := Run()
	require.NotNil(t, err, err)
}
