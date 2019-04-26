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
	"os"
	"os/signal"

	"github.com/untillpro/godif"
	"github.com/untillpro/godif/iservices"
)

var signals chan os.Signal

// Run calls godif.ResolveAll(), starts all services and wait until Terminate() is called
// When Terminate() is called ctx is cancelled and all Stop's are called asynchronously
// # Events
func Run() error {

	Declare()
	godif.Require(&iservices.Start)
	godif.Require(&iservices.Stop)

	errs := godif.ResolveAll()
	defer godif.Reset()
	if len(errs) > 0 {
		return errs
	}

	ctx, cancel := context.WithCancel(context.Background())

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer iservices.Stop(ctx)

	var err error
	ctx, err = iservices.Start(ctx)
	if nil != err {
		cancel()
		return err
	}

	sig := <-signals
	log.Println("[services] Signal received:", sig)
	cancel()
	return nil
}

// Terminate running Run
func Terminate() {
	signals <- os.Interrupt
}
