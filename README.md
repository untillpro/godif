# godif

Go dependency injection for functions (and not only...)

# Usage Example

Imagine we have a functional interface to work with key-value database. Database has two methods - `Put` and `Get`. These methods works with `BucketDef`.

The following example shows how to:

1. Declare functional interface
2. Implement (provide) functions
3. Use functional interface

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

Package declares variables of certain types, these variablesmust be furher provided.

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

## 3. Use

```go
package service

var bucketService = &ikvdb.BucketDef{Key: "service"}

// Declare requires Put function and provides `bucketService` definition
func Declare() {
	godif.Require(&ikvdb.Put)
	godif.ProvideExtension(&ikvdb.BucketDefs, bucketService.Key, bucketService)
}

type ctxKey string

// CtxUserName denotes user name
var CtxUserName = ctxKey("UserName")

// Start something
func Start(ctx context.Context) {
	user := ctx.Value(CtxUserName)
	ikvdb.Put(ctx, bucketService, "startedTime", time.Now())
	ikvdb.Put(ctx, bucketService, "startedBy", user)
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
