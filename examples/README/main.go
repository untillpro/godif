package main

import (
	"context"
	"log"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/kvdb"
	"github.com/untillpro/godif/examples/README/service"
)

func main() {
	kvdb.Declare()
	service.Declare()

	errs := godif.ResolveAll()
	if len(errs) != 0 {
		// Non-assignalble Requirements
		// Unresolved dependencies
		// Multiple provisions
		// Incompatible types
		// Non-required provisions
		log.Panic(errs)
	}

	ctx := context.WithValue(context.Background(), service.CtxUserName, "Peter")
	service.Start(ctx)
}
