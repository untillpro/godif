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

func Declare() {
    godif.Provide(&ikvdb.Put, Put)
    godif.Provide(&ikvdb.Get, Get)
}

// Get implements ikvdb.Get
func Get(ctx context.Context, key interface{}) (value interface{}, ok bool) {
	return nil, false
}


```

## 3. Use Functions

```go
package usage

func Declare() {
    godif.Require(&decl.Func1)
    godif.Require(&decl.Func2)
}

func

```

## 4. Build App

```go
package main

func main(){
    prov.Declare()
    usage.Declare()

    errs := godif.ResolveAll()
    if len(errs) != 0{
        // Non-assignalble Requirements
        // Cyclic dependencies
        // Unresolved dependencies
        // Multiple implementations
        log.Panic(errs)
    }

    // All implementors of godif.InitFunc will be called
    // Dependency defines the order of init
    errs = godif.Init()
    defer godif.Finit()

    if len(errs) != 0{
        log.Panic(errs)
    } 

    // Do something
    ikvdb.Put("key1", "value1")
    v1, ok := ikvdb.Get("key1")

}

```
# Under the Hood

- All registration functions works with default instance of `godif.ContainerDeclaration`
- ResolveAll creates default instance of `godif.ContainerInstance` which is used to init/finit/start/stop