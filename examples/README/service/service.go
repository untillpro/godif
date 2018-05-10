package service

import (
	"context"

	"github.com/maxim-ge/godif/examples/README/godif"
	"github.com/maxim-ge/godif/examples/README/ikvdb"
)

type s struct{}

// Declare dependencies and provisions
func Declare(cd *godif.CtxDecl) {
	cd.Require(&ikvdb.Get)
	cd.Require(&ikvdb.Put)
	//	cd.RegisterCtxInit(CtxInit)
	cd.RegisterCtxMain(CtxMain)
}

// CtxInit -ilalizer
func CtxInit(ctx context.Context, folderName string) context.Context {
	return ctx
}

// CtxMain function
func CtxMain(ctx context.Context) {
}
