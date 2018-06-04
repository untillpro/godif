/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"bytes"
	"fmt"
	"reflect"
)

// Errors is multiple errors container
type Errors []error

// EMultipleImplementations occurs if there are more than one provided implementation for one type
type EMultipleImplementations struct {
	Type  reflect.Type
	impls []srcElem
}

// EImplementationNotProvided error occurs if there is no implementation provided for a type
type EImplementationNotProvided struct {
	Type reflect.Type
}

// ENonAssignableRequirement error occurs if non-assignable (e.g. not variable) requirement is declared
type ENonAssignableRequirement struct {
	Type reflect.Type
	requirement srcElem
}

func (e *EMultipleImplementations) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.impls {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("multiple implementations of %s provided at:\r\n%s", e.Type, buffer.String())
}

func (e *EImplementationNotProvided) Error() string {
	return fmt.Sprintf("implementation of %s is not provided", e.Type)
}

func (e *ENonAssignableRequirement) Error() string {
	return fmt.Sprintf("non-assignable requirement for %s is declared at %s:%d", e.Type, e.requirement.file, e.requirement.line)
}

func (e Errors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	msg := "multiple errors:"
	for _, err := range e {
		msg += "\n" + err.Error()
	}
	return msg
}
