/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import "context"

// Services elements should be provided by others
var Services []IService

// InitAndStart all services
// First all Init methods are called in order of registration, then all Start methods
// StopAndFinit must be always called aftewards
var InitAndStart func(ctx context.Context) (newCtx context.Context, err error)

// StopAndFinit stops and finits services
// All Stops are called asynchronously
// When all Stop's finish, Finits are called in reverse order of their provisions
var StopAndFinit func(ctx context.Context)
