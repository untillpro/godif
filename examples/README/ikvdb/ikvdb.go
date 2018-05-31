package ikvdb

import "context"

// Put saves given key and value to some persistent storage
var Put func(ctx context.Context, key interface{}, value interface{})

// Get gets the value from some persistent storage
var Get func(ctx context.Context, key interface{}) (value interface{}, ok bool)
