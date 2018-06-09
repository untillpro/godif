package service

import (
	"context"
	"time"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
)

// Declare dependencies and provisions
func Declare() {
	godif.Require(&ikvdb.Put)
}

type ctxKey string

// CtxUserName denotes user name
var CtxUserName = ctxKey("UserName")

// Start something
func Start(ctx context.Context) {
	user := ctx.Value(CtxUserName)
	ikvdb.Put(ctx, "startedTime", time.Now())
	ikvdb.Put(ctx, "startedBy", user)
}
