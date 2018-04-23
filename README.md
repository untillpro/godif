# godif

Inject Functions

# Usage

## 1. Declare Functions

```go
package decl

// Declare function types explicitly, otherwise functions are matched by signature
type Func1Type func(x int, y int)
type Func2Type func(s string)

var Func1 Func1Type
var Func2 Func2Type
```

## 2. Provide Functions

```go
package prov

func Declare() {
    godif.Provide(decl.func1, MyFunc1)
    godif.Provide(decl.func2, MyFunc2)
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
package req

func Declare() {
    godif.Require(&decl.Func1)
    godif.Require(&decl.Func2)
}
```

## 4. Build App


```go
package main

func main(){
    prov.Declare()
    req.Declare()

    errs := godif.ResolveAll()
    if len(errs) != 0{
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
    declare.Func1(1, 2)
    declare.Func2("Hello")

}

```
# Under the Hood

- All registration functions works with default instance of `godif.ContainerDeclaration`
- ResolveAll creates default instance of `godif.ContainerInstance` which is used to init/finit/start/stop
