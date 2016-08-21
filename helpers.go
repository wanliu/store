package store

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

func MergeStruct(obj interface{}, m map[string]interface{}) error {
	s := structs.New(obj)

	for k, v := range m {
		f, ok := s.FieldOk(k)

		if !ok {
			return errors.New(fmt.Sprintf("object missing this field: %s", k))
		}

		if reflect.ValueOf(v).Kind() != f.Kind() {
			return errors.New(fmt.Sprint("map field type is mistake object field type by :%s", k))
		}

		f.Set(v)
	}
	return nil
}
