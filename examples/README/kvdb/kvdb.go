package kvdb

import (
	"context"

	"github.com/maxim-ge/godif/examples/README/godif"
	"github.com/maxim-ge/godif/examples/README/ikvdb"
)

// Declare requirements and provisions
func Declare(cd *godif.CtxDecl) {
	cd.Provide(&ikvdb.Get, Get)
}

// Get implements ikvdb.Get
func Get(ctx context.Context, key interface{}) (value interface{}, ok bool) {
	return nil, false
}
