# godif

Inject functions

# Usage

## 1. Declare functions

```go
package decl

type func1 func(x int, y int)
type func2 func(s string)
```

## 2. Implement functions

```go
package implementor

func Declare() {
    godif.Provide((decl.func1)(nil), MyFunc1)
    godif.Provide((decl.func2)(nil), MyFunc2)
}

func MyFunc1(x int, y int) {
...
}

func MyFunc2(s string) {
...
}

```

## 3. Require Functions

```go
package consumer

var F1 decl.Func1
var F2 decl.Func2

func Declare() {
    godif.Require(&f1)
    godif.Require(&f2)
}
```

## 4. Build App


```go
package main

func main(){
    implementor.Declare()
    consumer.Declare()

    errs := godif.ResolveAll()
    if len(errs) != 0{
        // Cyclic dependencies
        // Unresolved dependencies
        // Multiple implementations
        log.Panic(errs)
    }

    // All implementors of godif.InitFunc will be called
    errs = godif.Init()
    defer godif.Finit()

    if len(errs) != 0{
        log.Panic(errs)
    } 

    // Do something
    consumer.F1(1, 2)
    consumer.F2("Hello")

}

```

# Under the Hood

- All registration functions works with default instance of `godif.ContainerDeclaration`
- ResolveAll creates default instance of `godif.ContainerInstance` which is used to init/finit/start/stop
