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
 )

var signals chan os.Signal

// Run starts all services and wait until Terminate() is called
// All Start  methods are called in order of registration
// If any error during start occurs it is immediately returned
// When Terminate() is called ctx is cancelled and all Stop's are called asynchronously
// # Events
func Run() error {
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
	
	<-signals
	log.Println("[services] Termination signal received")
	cancel()
	return nil
}

// Terminate running Run
func Terminate(){
	signals<-os.Interrupt
}