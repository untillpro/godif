package icounters

import "context"

// Inc increments the counter
var Inc func(ctx context.Context, counterName string)
