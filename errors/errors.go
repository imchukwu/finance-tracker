package errors

import "fmt"

type StorageError struct {
    Op  string
    Err error
}

func (e *StorageError) Error() string {
    return fmt.Sprintf("storage error during %s: %v", e.Op, e.Err)
}

type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in %s: %s", e.Field, e.Msg)
}