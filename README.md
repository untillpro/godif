# godif

Go dependency injection for functions (and not only...)

# Usage

## Provide key-value

- Declare: `var MyMap map[string]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string]int{}`
    - Use Provide() -> error
  - Provide implementation: `godif.Provide(&MyMap, map[string]int{})`
    - Multiple implementations -> error
    - Non-map or map of incompatible key or value type -> error
  - No implementation -> error
- Provide data: `godif.ProvideKeyValue(&MyMap, "key1", 1)`
  - Multiple values per key -> error
  - Key or value of different types provided -> error

## Provide key-slice

- Declare: `var MyMap map[string][]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string][]int{}`
    - Use Provide() -> error
  - Provide implementation: `godif.Provide(&MyMap, map[string][]int{})`
    - Multiple implementations -> error
    - slice of incompatible element type -> error
  - No implementation -> error
- Add initial data if needed: `MyMap["key1"] = append(MyMap["key1"], 42)`
  - Further `godif.ProvideKeyValue()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideKeyValue(&MyMap, "key1", 1)`
  - `godif.ProvideKeyValue(&MyMap, "key1", 2)`
  - `godif.ProvideKeyValue(&MyMap, "key1", []int{3, 4})`

## Provide slice element

- Declare: `var MySlice []string`
- Implement: `godif.Provide(&MySlice, []string{})`
- Implement
  - Manually: `MySlice = []string{}`
    - `godif.Provide(&MySlice, []string{})` -> Implementation provided for non-nil error
  - Provide implementation: `godif.Provide(&MySlice, []string{})`
  - No implementation -> Implementation not provided error
- Add initial data if needed: `MySlice = append(MySlice, 42)`
  - Further `godif.ProvideSliceElement()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideSliceElement(&MySlice, "str1")`
  - `godif.ProvideSliceElement(&MySlice, []string{"str3", "str4"})`

## Reset all injections
- `godif.Reset()`
- Provided vars will be nilled
- Manually inited vars will be kept
- Data injected into manually inited vars will be kept

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
