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

func Register() {
    godif.RegisterImpl((decl.func1)(nil), MyFunc1)
    godif.RegisterImpl((decl.func2)(nil), MyFunc2)
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

var f1 decl.Func1
var f2 decl.Func2

func Register() {
    godif.RegisterDep(&f1)
    godif.RegisterDep(&f2)
}
```

## 4. Build App


```go
package main

func main(){
    implementor.Register()
    consumer.Register()

    errs := godif.ResolveAll()
    if len(errs) != 0{
        // Cyclic dependencies
        // Unresolved dependencies
        log.Panic(errs)
    }

    // All implementors of godif.InitFunc will be called
    errs = godif.Init()
    defer godif.Finit()

    if len(errs) != 0{
        log.Panic(errs)
    } 

    errs = godif.Start()
    if len(errs) != 0{
        log.Panic(errs)
    } 
    defer godif.Stop()

    // Do something

}

```

# Under the Hood

- All functions works using static instance of `godif.Container`
