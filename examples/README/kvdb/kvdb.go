package kvdb

import (
	"context"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
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

}

// Get implements ikvdb.Get
func Get(ctx context.Context, key interface{}) (value interface{}, ok bool) {
	return nil, false
}

//
func Init(parentCtx context.Context) context.Context {
	return nil
}
