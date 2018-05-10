package main

import (
	"log"

	"github.com/maxim-ge/godif/examples/README/godif"
	"github.com/maxim-ge/godif/examples/README/kvdb"
	"github.com/maxim-ge/godif/examples/README/service"
)

func main() {
	kvdb.Declare(godif.RootCD)
	service.Declare(godif.RootCD)

	errs := godif.ResolveAll()
	if len(errs) != 0 {
		// Non-assignalble Requirements
		// Cyclic dependencies
		// Unresolved dependencies
		// Multiple provisions
		log.Panic(errs)
	}

	// All implementors of godif.InitFunc will be called
	// Dependency defines the order of init
	// ctx, errs := godif.Init()
	// defer godif.Finit()

	if len(errs) != 0 {
		log.Panic(errs)
	}

}
