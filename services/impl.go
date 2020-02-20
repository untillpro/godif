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
	"runtime/debug"

	"github.com/untillpro/godif"
)

// Services should be provided by godif.ProvideSliceElement(&services.Services, ...)
var Services []IService

// SetVerbose changes logging defaults (by default verbose is true)
func SetVerbose(value bool) (prev bool) {
	prev = verboseOutput
	verboseOutput = value
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

// StartServices starts all services registered in Services
// Calls Services' Start methods in order of provision
// If any error/panic occurs it immediately returns
func StartServices(ctx context.Context) (newCtx context.Context, err error) {
	newCtx, started, err = Start(ctx, Services, verboseOutput)
	return newCtx, err
}

// StopServices calls all Stop methods of started services in reversed order of provision
func StopServices(ctx context.Context) {
	Stop(ctx, started, verboseOutput)
	started = []IService{}
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
// If service panics EPanic is returned as err
func Start(startingCtx context.Context, servicesToStart []IService, verbose bool) (newCtx context.Context, startedServices []IService, err error) {
	if verbose {
		logln("Starting services...")
	}
	newCtx = startingCtx
	var startingService IService
	defer func() {
		if r := recover(); r != nil {
			if verbose {
				logln(fmt.Sprintf("Service paniced: %v: %v\n%s", reflect.TypeOf(startingService), r, string(debug.Stack())))
			}
			err = &EPanic{PanicData: r, PanicedService: startingService}
		}
	}()
	for _, startingService = range servicesToStart {
		serviceName := reflect.TypeOf(startingService).String()
		if verbose {
			logln("Starting " + serviceName + "...")
		}
		newCtx, err = startingService.Start(newCtx)
		if nil != err {
			logln("Error starting service:", err)
			return
		}
		startedServices = append(startedServices, startingService)
	}
	if verbose {
		logln("All services started")
	}
	return
}

// Stop all services in given context
// Services must NOT panic
func Stop(ctx context.Context, startedServices []IService, verbose bool) {
	if verbose {
		logln("Stopping...")
	}
	for i := len(startedServices) - 1; i >= 0; i-- {
		service := startedServices[i]
		serviceName := reflect.TypeOf(service).String()
		if verbose {
			logln("Stopping " + serviceName + "...")
		}
		service.Stop(ctx)
	}
	if verbose {
		logln("All services stopped")
	}
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
	logln("Resolving dependencies...")
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
	if !verboseOutput {
		return
	}
	pargs := []interface{}{"[services]"}
	pargs = append(pargs, args...)

	log.Println(pargs...)
}

var verboseOutput = true
var started []IService
var signals chan os.Signal
