/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"context"
	"sync"
	"os"
	"log"

	isvc "github.com/untillpro/godif/iservices"
)

var started = []isvc.IService{}
var signals chan os.Signal

func implStart(ctx context.Context) (context.Context, error) {

	log.Println("[services] Starting...")
	for _, service := range isvc.Services {
		var err error
		ctx, err = service.Start(ctx)
		if nil != err{
			log.Println("[services] Error starting services", err)
			return ctx, err
		}
		started = append(started, service)
	}
	log.Println("[services] Started")
	return ctx, nil
}

func implStop(ctx context.Context) {
	log.Println("[services] Stopping...")
	var wg sync.WaitGroup
	for _, service := range started {
		s := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop(ctx)
		}()
	}
	wg.Wait()
	started = []isvc.IService{}
	log.Println("[services] All services stopped")
}
