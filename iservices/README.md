# Overview

Declare, init, start, stop and finit services.

```go
import (
    "github.com/untillpro/godif/iservices"
    "github.com/untillpro/godif/services"    
)

    // We will need InitAndStart and StopAndFinit
    godif.Require(&iservices.InitAndStart)
    godif.Require(&iservices.StopAndFinit)

    // Declare implementation
    services.Declare()

    // Provide services
    godif.ProvideSliceElement(&iservices.Services, ...)
    godif.ProvideSliceElement(&iservices.Services, ...)

    // Resole all
    errs := godif.ResolveAll()

    // Start services
    ctx, err := iservices.InitAndStart(ctx)

    // Use services
    ...

    // Stop services
    iservices.StopAndFinit(ctx)

```    
- [Interface declaration](interface.go)
- [Interface test](interfacetest.go)