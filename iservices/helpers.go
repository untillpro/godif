/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

 package iservices

 import (
	"context"
	"os"
	"os/signal"
	"log"
	"github.com/untillpro/godif"
 )

var signals chan os.Signal

// Run calls godif.ResolveAll(), starts all services and wait until Terminate() is called
// When Terminate() is called ctx is cancelled and all Stop's are called asynchronously
// # Events
func Run() error {

	errs := godif.ResolveAll()
	if len(errs) > 0{
		return errs
	}
	defer godif.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer Stop(ctx)

	var err error
	ctx, err = Start(ctx)
	if nil != err{
		cancel()
		return err
	}
	
	sig := <- signals
	log.Println("[services] Signal received:", sig)
	cancel()
	return nil
}

// Terminate running Run
func Terminate(){
	signals<-os.Interrupt
}