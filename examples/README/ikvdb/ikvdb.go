package ikvdb

import "context"

// Put saves given key and value to some persistent storage
var Put func(ctx context.Context, bucket *BucketDef, key interface{}, value interface{})

// Get gets the value from some persistent storage
var Get func(ctx context.Context, bucket *BucketDef, key interface{}) (value interface{}, ok bool)

// BucketDef defines the bucket
type BucketDef struct {
	Key string
}

// BucketDefs keeps list of BucketDef
var BucketDefs map[string]*BucketDef
