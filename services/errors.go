/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package services

import "fmt"

// EPanic is returned by Start/Stop if some service paniced
type EPanic struct {
	service   IService
	panicData interface{}
}

func (e *EPanic) Error() string {
	return fmt.Sprintf("Service %v paniced: %v", e.service, e.panicData)
}
