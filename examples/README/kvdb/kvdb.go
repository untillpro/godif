package kvdb

import (
	"context"
	"log"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
)

// Declare provides Get/Put functions and map of BucketDef
func Declare() {
	godif.Provide(&ikvdb.Get, Get)
	godif.Provide(&ikvdb.Put, Put)
	godif.Provide(&ikvdb.BucketDefs, map[string]*ikvdb.BucketDef{})
}

var buckets = map[string]map[interface{}]interface{}{}

// Get implements ikvdb.Get
func Get(ctx context.Context, bucket *ikvdb.BucketDef, key interface{}) (value interface{}, ok bool) {
	kv, ok := buckets[bucket.Key]
	if !ok {
		log.Panicln("Bucket not found", bucket.Key)
	}
	val, ok := kv[key]
	return val, ok
}

// Put implements ikvdb.Put
func Put(ctx context.Context, bucket *ikvdb.BucketDef, key interface{}, value interface{}) {
	kv, ok := buckets[bucket.Key]
	if !ok {
		log.Panicln("Bucket not found", bucket.Key)
	}
	kv[key] = value
}
