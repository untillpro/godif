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
	req   *srcElem
	provs []*srcElem
}

// EImplementationNotProvided error occurs if there is no implementation provided for a type
type EImplementationNotProvided struct {
	req *srcElem
}

// EImplementationProvidedForNonNil error occurs if target value is not nil but implementation provided
type EImplementationProvidedForNonNil struct {
	req *srcElem
}

// ENonAssignableRequirement error occurs if non-assignable (e.g. not variable) requirement is declared
type ENonAssignableRequirement struct {
	req *srcElem
}

// EIncompatibleTypesFunc error occurs if type of a requirement (func) is incompatible to provided implementation
type EIncompatibleTypesFunc struct {
	req  *srcElem
	prov *srcElem
}

// EIncompatibleTypesSlice error occurs if type of a requirement (array or slice) is incompatible to provided implementation
type EIncompatibleTypesSlice struct {
	reqType reflect.Type
	prov    *srcElem
}

// EProvidedNotUsed error occurs if something is provided but not required
type EProvidedNotUsed struct {
	prov *srcElem
}

// EMultipleValues error occurs if more than one value is provided per one key by ProvideMapValue() call
type EMultipleValues struct {
	provs []*srcElem
}

// EAlreadyResolved occurs on ResolveAll() call if it called already
type EAlreadyResolved struct {
	resolvePlace *src
}

func (e *EMultipleImplementations) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.provs {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("Requirement at %s:%d has multiple provisions at:\r\n%s", e.req.file, e.req.line, buffer.String())
}

func (e *EImplementationNotProvided) Error() string {
	return fmt.Sprintf("Requirement %T at %s:%d is not provided", e.req.elem, e.req.file, e.req.line)
}

func (e *EImplementationProvidedForNonNil) Error() string {
	return fmt.Sprintf("Implementation provided for non-nil %T at %s:%d", e.req.elem, e.req.file, e.req.line)
}

func (e *ENonAssignableRequirement) Error() string {
	return fmt.Sprintf("Non-assignable requirement at %s:%d", e.req.file, e.req.line)
}

func (e *EIncompatibleTypesSlice) Error() string {
	return fmt.Sprintf("Incompatible types: %s required but %s provided at %s:%d", e.reqType,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EIncompatibleTypesFunc) Error() string {
	return fmt.Sprintf("Incompatible types: %s required at %s:%d, %s provided at %s:%d", reflect.TypeOf(e.req.elem), e.req.file, e.req.line,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EProvidedNotUsed) Error() string {
	return fmt.Sprintf("%T Provided at %s:%d but not used", e.prov.elem, e.prov.file, e.prov.line)
}

func (e *EMultipleValues) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.provs {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("Extension point has multiple values provided at:\r\n%s", buffer.String())
}

func (e *EAlreadyResolved) Error() string {
	return fmt.Sprintf("Already resolved at %s:%d", e.resolvePlace.file, e.resolvePlace.line)
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
