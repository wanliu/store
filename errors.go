package store

import (
	"errors"
	"fmt"
)

var (
	ErrNonNamedField      = errors.New("can't found this named field of struct")
	ErrUnqiueAwayls       = errors.New("awayls have unique index or index value is empty")
	ErrInvalidPKey        = errors.New("primary key is invalid number")
	ErrMissingBucket      = errors.New("cant open this bucket")
	ErrMissingIndexBucket = errors.New("can't open this index bucket")
	ErrZeroID             = errors.New("ID can't be zero value")
	ErrNotEmptyIndex      = errors.New("index is not empty")
	ErrConvertInvalid     = errors.New("invalid type convert")
	ErrMissingValue       = errors.New("value is nil or empty")
)

type ErrorUniqueIndex struct {
	store Store
	index Index
	value interface{}
}
type ErrorInvalidField struct {
	obj  interface{}
	name string
}

func (err ErrorUniqueIndex) Error() string {
	return fmt.Sprintf("store: %s awayls have same index named `%s` equals `%v` or variable is empty", err.store, err.value, err.index)
}

func (err ErrorInvalidField) Error() string {
	return fmt.Sprintf("store: %s can't found field `%s` in %v struct", err.name, err.obj)
}
