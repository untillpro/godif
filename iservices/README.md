# Overview

Declare, init, start, stop and finit services. Ref. [Test_BasicUsage](service_test.go)

```go
import (
    "github.com/untillpro/godif/iservices"
)

    // We will need InitAndStart and StopAndFinit
    godif.Require(&iservices.InitAndStart)
    godif.Require(&iservices.StopAndFinit)

    // Declare default implementation
    iservices.Declare()

    // Provide services
    godif.ProvideSliceElement(&iservices.Services, ...)
    godif.ProvideSliceElement(&iservices.Services, ...)

    // Resole all
    errs := godif.ResolveAll()

    // Start services
    ctx, err := iservices.InitAndStart(ctx)

    // Stop services
    iservices.StopAndFinit(ctx)

```    