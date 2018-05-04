package counters

import (
	"context"

	"github.com/maxim-ge/godif/examples/1/api/icounters"
	"github.com/maxim-ge/godif/examples/1/godif"
)

// Declare requirements and provisions
func Declare() {
	godif.Provide(icounters.Inc, Inc)
}

// Inc counter with given name
func Inc(ctx context.Context, counterName string) {

}
