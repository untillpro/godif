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

    if justRun {
        // ResoleAll(), starts all services and wait until os.Interrupt received
        // Can be terminated also using iservices.Terminate()
        err := iservices.Run()
    } else {

        // Resole all
        errs := godif.ResolveAll()

        if len(errs) >0 {
            log.Fatal("Resolve error", errs)
        }
        
        defer godif.Reset()

        // Start services
        ctx, err := iservices.Start(context.Background())
        defer iservices.Stop(ctx)

        // Use services
        ...
    }

```    
- [Interface declaration](interface.go)
- [Interface test](interfacetest.go)