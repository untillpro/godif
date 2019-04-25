/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import(
	"context"
)

// Services elements should be provided by others
var Services []IService

// Start calls Start methods in order of provision
// If any error it is immediately returned
// Ref also Run() in helpers.go
var Start func(ctx context.Context) (context.Context, error)

// Stop asyncronously calls all Stop methods
var Stop func(ctx context.Context)
