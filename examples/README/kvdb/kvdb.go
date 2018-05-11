package kvdb

import (
	"context"

	"github.com/maxim-ge/godif/examples/README/godif"
	"github.com/maxim-ge/godif/examples/README/ikvdb"
)

// Params are used for CtxInst initialization
type Params struct {
	Dir      string
	ValueDir string
}

// DeclareDb requirements and provisions
func DeclareDb(cd *godif.CtxDecl) {
	cd.Provide(&ikvdb.Get, Get)
}

// DeclareTransaction requirements and provisions
func DeclareTransaction(cd *godif.CtxDecl) {
	cd.Provide
	cd.
}

// Get implements ikvdb.Get
func Get(ctx context.Context, key interface{}) (value interface{}, ok bool) {
	return nil, false
}

//
func Init(parentCtx context.Context) context.Context {

}
