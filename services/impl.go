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
	"reflect"

	"github.com/untillpro/godif"
)

// Services should be provided by godif.ProvideSliceElement(&services.Services, ...)
var Services = []IService{}

var started []IService

// DisableLogging s.e.
// By default logging is on
func DisableLogging() {
	loggingEnabled = false
}

// StartServices starts all services
// Calls Services' Start methods in order of provision
// If any error occurs it is immediately returned
func StartServices(ctx context.Context) (context.Context, error) {

	logln("Starting services...")
	for _, service := range Services {
		var err error
		serviceName := reflect.TypeOf(service).String()
		logln("Starting " + serviceName + "...")
		ctx, err = service.Start(ctx)
		if nil != err {
			logln("Error starting service:", err)
			return ctx, err
		}
		started = append(started, service)
	}
	logln("All services started")
	return ctx, nil
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

// Run calls godif.ResolveAll(), starts all services and wait until Terminate() is called
// When Terminate() is called ctx is cancelled and all Stop's are called asynchronously
// # Events
func Run() error {

	errs := godif.ResolveAll()
	defer godif.Reset()
	if len(errs) > 0 {
		return errs
	}

	ctx, cancel := context.WithCancel(context.Background())

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer StopServices(ctx)

	var err error
	ctx, err = StartServices(ctx)
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

// ResolveAndStart resolve deps and starts services
func ResolveAndStart() (context.Context, error) {
	err := godif.ResolveAll()
	if nil != err {
		return context.Background(), err
	}
	return StartServices(context.Background())
}

// StopAndReset stops services and resets deps
func StopAndReset(ctx context.Context) {
	StopServices(ctx)
	godif.Reset()
}
