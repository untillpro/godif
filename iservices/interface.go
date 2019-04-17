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
// If any error occurs it is immediately returned
// StopAndFinit must be always called aftewards
var InitAndStart func(ctx context.Context) (newCtx context.Context, err error)

// StopAndFinit stops and finits services
// All Stop's are called asynchronously
// Stop's called only if appropriate Start() succeeded
// When all Stop's finish, Finit's are called in reverse order of their provisions
// Finit is called only if Init succeeded
var StopAndFinit func(ctx context.Context)
