/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"github.com/untillpro/godif"
	"github.com/untillpro/godif/iservices"
)

// Declare requirements and provisions
func Declare() {
	godif.Provide(&iservices.Start, implStart)
	godif.Provide(&iservices.Stop, implStop)
	godif.Provide(&iservices.Services, make([]iservices.IService, 0, 50))
}
