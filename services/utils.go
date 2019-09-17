/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"context"
	"log"
	"reflect"
)

var verboseEnabled = true

// Start given services in given context using IService.Start() method
// If IService.Start() returns error function exists
// newCtx: new context
// started: services which were succedsfully started
// err: error reported by service
func Start(startingCtx context.Context, servicesToStart []IService) (newCtx context.Context, startedServices []IService, err error) {
	logln("Starting services...")
	newCtx = startingCtx
	for _, service := range servicesToStart {
		serviceName := reflect.TypeOf(service).String()
		logln("Starting " + serviceName + "...")
		newCtx, err = service.Start(newCtx)
		if nil != err {
			logln("Error starting service:", err)
			return
		}
		startedServices = append(startedServices, service)
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

func logln(args ...interface{}) {
	if !verboseEnabled {
		return
	}
	pargs := []interface{}{"[services]"}
	pargs = append(pargs, args...)

	log.Println(pargs...)
}
