/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"context"

	"github.com/untillpro/godif"
)

// ResolveAndStart resolve deps and starts services
func ResolveAndStart() (context.Context, error) {
	err := godif.ResolveAll()
	if nil != err {
		return context.Background(), err
	}
	return StartServices(context.Background())
}

// StopAndReset stops services and resets deps
func StopAndReset(ctx context.Context) {
	StopServices(ctx)
	godif.Reset()
}
