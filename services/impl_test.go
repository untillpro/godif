/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import (
	"testing"

	"github.com/untillpro/godif/iservices"
)

func Test_Impl(t *testing.T) {
	iservices.TestImpl(t, Declare)
}
