# godif

Go dependency injection for functions (and not only...)

# Usage

- Package `ikvdb` declares functional interface (`Put`, `Get`) to buckets and  bucket definitions holder (`Buckets`)
- Package `kvdb` provides functions and bucket definitions holder

## 1. Declare

```go
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

```

## 2. Provide

```go

package kvdb

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
```

## 3. Use Functions

```go
package service

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

```

## 4. Build App

```go
func main() {
	kvdb.Declare()
	service.Declare()

	errs := godif.ResolveAll()
	if len(errs) != 0 {
		// Non-assignalble Requirements
		// Unresolved dependencies
		// Multiple provisions
		log.Panic(errs)
	}

	ctx := context.WithValue(context.Background(), service.CtxUserName, "Peter")
	service.Start(ctx)
}
```

# Declare/Provide Data

