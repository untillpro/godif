# godif

Go dependency injection for functions (and not only...)

# Usage

## 1. Declare Functions

```go
package ikvdb

// Put saves given key and value to some persistent storage
var Put func(ctx context.Context, key interface{}, value interface{})

// Get gets the value from some persistent storage
var Get func(ctx context.Context, key interface{}) (value interface{}, ok bool)

```

## 2. Provide Functions

```go
package kvdb

// Declare requirements and provisions
func Declare() {
	godif.Provide(&ikvdb.Get, Get)
	godif.Provide(&ikvdb.Put, Put)
}

var mapDb = make(map[interface{}]interface{})

// Get implements ikvdb.Get
func Get(ctx context.Context, key interface{}) (value interface{}, ok bool) {
	val, ok := mapDb[key]
	return val, ok
}

// Put implements ikvdb.Put
func Put(ctx context.Context, key interface{}, value interface{}) {
	mapDb[key] = value
}
```

## 3. Use Functions

```go
package service

// Declare dependencies and provisions
func Declare() {
	godif.Require(&ikvdb.Put)
}

type ctxKey string

// CtxUserName denotes user name
var CtxUserName = ctxKey("UserName")

// Start something
func Start(ctx context.Context) {
	user := ctx.Value("CurrentUser")
	ikvdb.Put(ctx, "startedTime", time.Now())
	ikvdb.Put(ctx, "startedBy", user)
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
		// Cyclic dependencies
		// Unresolved dependencies
		// Multiple provisions
		log.Panic(errs)
	}

	ctx := context.WithValue(context.Background(), service.CtxUserName, "Peter")
	service.Start(ctx)
}
```
# Under the Hood

- All registration functions works with default instance of `godif.ContainerDeclaration`
- ResolveAll creates default instance of `godif.ContainerInstance` which is used to init/finit/start/stop