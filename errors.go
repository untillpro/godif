package godif

import (
	"fmt"
	"reflect"
)

// Errors is multiple errors container
type Errors []error

// EMultipleImplementations occurs if there are more than one provided implementation for one type
type EMultipleImplementations struct {
	Type reflect.Type
	Count int
	// sources slice (line numbers etc)
}

// EImplementationNotProvided error occurs if there is no implementation provided for a type
type EImplementationNotProvided struct {
	Type reflect.Type
	// sources slice (line numbers etc), check Type - if has already
}

// ENonAssignableRequirement error occurs if non-assignable (e.g. not variable) requirement is declared
type ENonAssignableRequirement struct {
	Type reflect.Type
}

func (e *EMultipleImplementations) Error() string {
	return fmt.Sprintf("%s implementation provided %d times", e.Type, e.Count)
}

func (e *EImplementationNotProvided) Error() string {
	return fmt.Sprintf("implementation of %s is not provided", e.Type)
}

func (e *ENonAssignableRequirement) Error() string {
	return fmt.Sprintf("non-assignalble requirement for %s is declared", e.Type)
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