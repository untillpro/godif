/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
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

// Services should be provided by godif.ProvideSliceElement(&services.Services, ...)
var Services = []IService{}

var started []IService

// StartServices starts all services
// Calls Services' Start methods in order of provision
// If any error occurs it is immediately returned
func StartServices(ctx context.Context) (context.Context, error) {

	log.Println("[services] Starting services...")
	for _, service := range Services {
		var err error
		serviceName := reflect.TypeOf(service).String()
		log.Println("[services] Starting " + serviceName + "...")
		ctx, err = service.Start(ctx)
		if nil != err {
			log.Println("[services] Error starting service:", err)
			return ctx, err
		}
		started = append(started, service)
	}
	log.Println("[services] All services started")
	return ctx, nil
}

// StopServices calls all Stop methods of started services in reversed order of provision
func StopServices(ctx context.Context) {
	log.Println("[services] Stopping...")
	for i := len(started) - 1; i >= 0; i-- {
		s := started[i]
		s.Stop(ctx)
	}
	started = []IService{}
	log.Println("[services] All services stopped")
}
