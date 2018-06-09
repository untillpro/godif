package service

import (
	"context"
	"time"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
)

var bucketStart = &ikvdb.BucketDef{Key: "start"}

// Declare requires Put function and provides `bucketStart` definition
func Declare() {
	godif.Require(&ikvdb.Put)
	godif.ProvideMapValue(&ikvdb.BucketDefs, bucketStart)
}

type ctxKey string

// CtxUserName denotes user name
var CtxUserName = ctxKey("UserName")

// Start something
func Start(ctx context.Context) {
	user := ctx.Value(CtxUserName)
	ikvdb.Put(ctx, bucketStart, "startedTime", time.Now())
	ikvdb.Put(ctx, bucketStart, "startedBy", user)
}
