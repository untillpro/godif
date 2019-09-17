/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"context"
	"os"
	"os/signal"
	"github.com/untillpro/godif"
)

// Services should be provided by godif.ProvideSliceElement(&services.Services, ...)
var Services []IService

var started []IService

// SetVerbose changes logging defaults (by default verbose is true)
func SetVerbose(value bool) (prev bool) {
	prev = verboseEnabled
	verboseEnabled = value
	return
}

// StartServices starts all services
// Calls Services' Start methods in order of provision
// If any error occurs it is immediately returned
func StartServices(ctx context.Context) (newCtx context.Context, err error) {
	newCtx, started, err = Start(ctx, Services)
	return newCtx, err
}

// StopServices calls all Stop methods of started services in reversed order of provision
func StopServices(ctx context.Context) {
	logln("Stopping...")
	for i := len(started) - 1; i >= 0; i-- {
		s := started[i]
		s.Stop(ctx)
	}
	started = []IService{}
	logln("All services stopped")
}

var signals chan os.Signal

// Declare s.e.
func Declare() {
	godif.Provide(&Services, []IService{})
}

// Run calls godif.ResolveAll(), starts all services and wait until Terminate() is called
// When Terminate() is called ctx is cancelled and all Stop's are called asynchronously
// # Events
func Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ctx, err := ResolveAndStartCtx(ctx)
	defer StopAndReset(ctx)

	if nil != err {
		cancel()
		return err
	}

	sig := <-signals
	logln("Signal received:", sig)
	cancel()
	return nil
}

// Terminate running Run
func Terminate() {
	signals <- os.Interrupt
}

// ResolveAndStart calls ResolveAndStartCtx  with background context
func ResolveAndStart() (context.Context, error) {
	return ResolveAndStartCtx(context.Background())
}

// ResolveAndStartCtx declares own provisions, resolve deps and starts services with given context
func ResolveAndStartCtx(ctx context.Context) (context.Context, error) {
	Declare()
	err := godif.ResolveAll()
	if nil != err {
		return ctx, err
	}
	return StartServices(ctx)
}

// StopAndReset stops services and resets deps
func StopAndReset(ctx context.Context) {
	StopServices(ctx)
	godif.Reset()
}
