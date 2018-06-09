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
	req   srcElem
	impls []srcElem
}

// EImplementationNotProvided error occurs if there is no implementation provided for a type
type EImplementationNotProvided struct {
	req srcElem
}

// ENonAssignableRequirement error occurs if non-assignable (e.g. not variable) requirement is declared
type ENonAssignableRequirement struct {
	req srcElem
}

// EIncompatibleTypes error occurs if type of a requirement is incompatible to provided implementation
type EIncompatibleTypes struct {
	req  srcElem
	impl srcElem
}

func (e *EMultipleImplementations) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.impls {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("Requirement at %s:%d has multiple provisions at:\r\n%s", e.req.file, e.req.line, buffer.String())
}

func (e *EImplementationNotProvided) Error() string {
	return fmt.Sprintf("Requirement at %s:%d is not provided", e.req.file, e.req.line)
}

func (e *ENonAssignableRequirement) Error() string {
	return fmt.Sprintf("Non-assignable requirement at %s:%d", e.req.file, e.req.line)
}

func (e *EIncompatibleTypes) Error() string {
	return fmt.Sprintf("Incompatible types: %s required at %s:%d, %s provided at %s:%d", reflect.TypeOf(e.req.elem), e.req.file, e.req.line,
		reflect.TypeOf(e.impl.elem), e.impl.file, e.impl.line)
}

func (e Errors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	msg := "Multiple errors:"
	for _, err := range e {
		msg += "\n" + err.Error()
	}
	return msg
}
