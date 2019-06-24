/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import "log"

var loggingEnabled = true

func logln(args ...interface{}) {
	if !loggingEnabled {
		return
	}
	pargs := []interface{}{"[services]"}
	pargs = append(pargs, args...)

	log.Println(pargs...)
}
