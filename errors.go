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
	"strings"
)

// Errors is multiple errors container
type Errors []error

// EMultipleStorageImplementations occurs if there are more than one implementations provided for slice or map
type EMultipleStorageImplementations struct {
	provs []*srcPkgElem
}

// EMultipleFuncImplementations occurs if there are more than one implementations provided for one func
type EMultipleFuncImplementations struct {
	req   *srcElem
	provs []*srcPkgElem
}

// EImplementationNotProvided error occurs if there is no implementation provided for a type
type EImplementationNotProvided struct {
	target *srcElem
}

// EImplementationProvidedForNonNil error occurs if target value is not nil but implementation provided
type EImplementationProvidedForNonNil struct {
	prov *srcPkgElem
}

// ENonAssignableRequirement error occurs if non-assignable (e.g. not variable) requirement is declared
type ENonAssignableRequirement struct {
	target *srcElem
}

// EIncompatibleTypesFunc error occurs if type of a requirement (func) is incompatible to provided implementation
type EIncompatibleTypesFunc struct {
	target *srcElem
	prov   *srcPkgElem
}

// EIncompatibleTypesStorageValue error occurs if type of an array or slice element or value of map is incompatible to provided implementation
type EIncompatibleTypesStorageValue struct {
	targetType reflect.Type
	prov       *srcElem
}

// EIncompatibleTypesStorageKey error occurs if type of key map is incompatible to provided implementation
type EIncompatibleTypesStorageKey struct {
	targetType reflect.Type
	prov       *srcElem
}

// EIncompatibleTypesStorageImpl error occurs if type of array or slice or map is incompatible to provided implementation
type EIncompatibleTypesStorageImpl struct {
	targetType reflect.Type
	prov       *srcElem
}

// EPackageNotUsed s.e.
type EPackageNotUsed struct {
	pkgName string
}

// EMultipleValues error occurs if more than one value is provided per one key by ProvideMapValue() call
type EMultipleValues struct {
	provs []*srcElem
}

// EAlreadyResolved occurs on ResolveAll() call if it called already
type EAlreadyResolved struct {
	resolvePlace *src
}

func (e *EMultipleStorageImplementations) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.provs {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("Multiple provisions of one storage at:\r\n%s", buffer.String())
}

func (e *EMultipleFuncImplementations) Error() string {
	var buffer bytes.Buffer
	for _, impl := range e.provs {
		buffer.WriteString(fmt.Sprintf("\t%s:%d\r\n", impl.file, impl.line))
	}

	return fmt.Sprintf("Requirement at %s:%d has multiple provisions at:\r\n%s", e.req.file, e.req.line, buffer.String())
}

func (e *EImplementationNotProvided) Error() string {
	kind := reflect.TypeOf(e.target.elem).Elem().Kind()
	if kind == reflect.Func {
		return fmt.Sprintf("Implementation of %T at %s:%d is not provided", e.target.elem, e.target.file, e.target.line)
	}
	return fmt.Sprintf("Target %T is nil at %s:%d. Init it manually or use Provide()", e.target.elem, e.target.file, e.target.line)
}

func (e *EImplementationProvidedForNonNil) Error() string {
	return fmt.Sprintf("Implementation provided for non-nil %T at %s:%d", e.prov.elem, e.prov.file, e.prov.line)
}

func (e *ENonAssignableRequirement) Error() string {
	return fmt.Sprintf("Non-assignable requirement at %s:%d. Use pointers to target on Require() and Provide()", e.target.file, e.target.line)
}

func (e *EIncompatibleTypesStorageValue) Error() string {
	return fmt.Sprintf("Incompatible types: target is %s but %s is used as value at %s:%d", e.targetType,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EIncompatibleTypesStorageKey) Error() string {
	return fmt.Sprintf("Incompatible types: target is %s but %s is used as key at %s:%d", e.targetType,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EIncompatibleTypesStorageImpl) Error() string {
	return fmt.Sprintf("Incompatible types: target is %s but %s is provided at %s:%d", e.targetType,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EIncompatibleTypesFunc) Error() string {
	return fmt.Sprintf("Incompatible types: %s required at %s:%d, %s provided at %s:%d", reflect.TypeOf(e.target.elem), e.target.file, e.target.line,
		reflect.TypeOf(e.prov.elem), e.prov.file, e.prov.line)
}

func (e *EPackageNotUsed) Error() string {
	return fmt.Sprintf("Have provisions from package %s but nothing is required from this package", e.pkgName)
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

	var sb strings.Builder
	sb.WriteString("Multiple errors:")
	for _, err := range e {
		sb.WriteString("\n" + err.Error())
	}
	return sb.String()
}
