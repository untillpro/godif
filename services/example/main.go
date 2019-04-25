/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

 package main

 import(
	"github.com/untillpro/godif"
	"github.com/untillpro/godif/iservices"
	"github.com/untillpro/godif/services"
	"context"
	"fmt"
 )

 func main(){

	// Register services

	godif.ProvideSliceElement(&iservices.Services, new(myService1))
	godif.ProvideSliceElement(&iservices.Services, new(myService2))

	// Will be terminated by SIGTERM
	services.Run()
 }

type myService1 struct{
}

func (s *myService1) Start(ctx context.Context) (context.Context, error){
	fmt.Println("Service1 started")
	return ctx, nil
 }

func  (s *myService1) Stop(ctx context.Context){
	fmt.Println("Service1 stopped")
}

type myService2 struct{
}

func (s *myService2) Start(ctx context.Context) (context.Context, error){
	fmt.Println("Service2 started")
	return ctx, nil
 }

func  (s *myService2) Stop(ctx context.Context){
	fmt.Println("Service2 stopped")
}

