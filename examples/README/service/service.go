package service

import (
	"context"
	"time"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
)

var bucketService = &ikvdb.BucketDef{Key: "service"}

// Declare requires Put function and provides `bucketStart` definition
func Declare() {
	godif.Require(&ikvdb.Put)
	godif.ProvideMapValue(&ikvdb.BucketDefs, bucketService.Key, bucketService)
}

type ctxKey string

// CtxUserName denotes user name
var CtxUserName = ctxKey("UserName")

// Start something
func Start(ctx context.Context) {
	user := ctx.Value(CtxUserName)
	ikvdb.Put(ctx, bucketService, "startedTime", time.Now())
	ikvdb.Put(ctx, bucketService, "startedBy", user)
}
