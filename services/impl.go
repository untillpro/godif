/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"

	"github.com/untillpro/godif"
)

// Services should be provided by godif.ProvideSliceElement(&services.Services, ...)
var Services []IService

// SetVerbose changes logging defaults (by default verbose is true)
func SetVerbose(value bool) (prev bool) {
	prev = verboseEnabled
	verboseEnabled = value
	return
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

// Declare s.e.
func Declare() {
	godif.Provide(&Services, []IService{})
}

// Start given services in given context using IService.Start() method
// If IService.Start() returns error function exists
// newCtx: new context
// started: services which were succedsfully started
// err: error reported by service
func Start(startingCtx context.Context, servicesToStart []IService) (newCtx context.Context, startedServices []IService, err error) {
	logln("Starting services...")
	newCtx = startingCtx
	var startingService IService
	defer func() {
		if r := recover(); r != nil {
			logln(fmt.Sprintf("Service paniced: %v: %v", startingService, r))
			err = &EPanic{PanicData: r, PanicedService: startingService}
		}
	}()
	for _, startingService = range servicesToStart {
		serviceName := reflect.TypeOf(startingService).String()
		logln("Starting " + serviceName + "...")
		newCtx, err = startingService.Start(newCtx)
		if nil != err {
			logln("Error starting service:", err)
			return
		}
		startedServices = append(startedServices, startingService)
	}
	logln("All services started")
	return
}

// Stop all services in given context
func Stop(ctx context.Context, startedServices []IService) {
	logln("Stopping...")
	for i := len(startedServices) - 1; i >= 0; i-- {
		service := startedServices[i]
		serviceName := reflect.TypeOf(service).String()
		logln("Stopping " + serviceName + "...")
		service.Stop(ctx)
	}
	logln("All services stopped")
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

func logln(args ...interface{}) {
	if !verboseEnabled {
		return
	}
	pargs := []interface{}{"[services]"}
	pargs = append(pargs, args...)

	log.Println(pargs...)
}

var verboseEnabled = true
var started []IService
var signals chan os.Signal
