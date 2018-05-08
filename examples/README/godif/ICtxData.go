package godif

import "context"

// ICtxData inits/finits some context ata
type ICtxData interface {
	Init(ctx context.Context) context.Context
	// panic is nil if no panic happens (CtxMain and previous Finit's)
	Finit(ctx context.Context, panic error)
}
