package kvdb

import (
	"context"

	"github.com/untillpro/godif/examples/README/godif"
	"github.com/untillpro/godif/examples/README/ikvdb"
)

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
