# Overview

Declare, init, start, stop and finit services.

```go
import (
    "github.com/untillpro/godif/iservices"
    "github.com/untillpro/godif/services"    
    "github.com/untillpro/godif"
)

    // Provice service implementation
    services.Declare()

    // Provide services
    godif.ProvideSliceElement(&iservices.Services, ...)
    godif.ProvideSliceElement(&iservices.Services, ...)

    // Resole all
    errs := godif.ResolveAll()

    if justRun {
        // Starts all services and wait until os.Interrupt received
        // Can be terminated also using iservices.Terminate()
        err := iservices.Run()
    } else {
        
        // Start services
        ctx, err := iservices.Start(context.Background())

        // Use services
        ...

        // Stop services
        iservices.Stop(ctx)

    }

```    
- [Interface declaration](interface.go)
- [Interface test](interfacetest.go)