/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import "github.com/untillpro/godif"

// Declare s.e.
func Declare() {
	godif.Provide(&InitAndStart, initAndStartImpl)
	godif.Provide(&StopAndFinit, stopAndFinitImpl)
	godif.Provide(&Services, make([]IService, 0, 50))
}
