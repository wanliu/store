package store

import (
	"bytes"
	"encoding/binary"

	"github.com/fatih/structs"
)

type primary_key struct {
	name string
	val  *structs.Field
}

type PkeyFunc func(Keyer) error

type Keyer interface {
	Set(uint64) error
	Value() uint64
	Name() string
}

func PrimaryKey(obj interface{}, feld string) (*primary_key, error) {
	s := structs.New(obj)
	val, ok := s.FieldOk(feld)

	if !ok {
		return nil, ErrNonNamedField
	}

	return &primary_key{val: val, name: feld}, nil
}

func (pkey *primary_key) Set(ID uint64) error {
	return pkey.val.Set(ID)
}

func (pkey *primary_key) Name() string {
	return pkey.name
}

func (pkey *primary_key) Value() uint64 {
	return pkey.val.Value().(uint64)
}

func (s *store) IdName(key Keyer) []byte {
	var (
		bval = make([]byte, 8)
		bkey = []byte(key.Name())
	)

	binary.LittleEndian.PutUint64(bval, key.Value())
	bs := [][]byte{bkey, bval}

	return bytes.Join(bs, bSep)
}

func (s *store) primaryKey(obj interface{}) (Keyer, error) {
	return PrimaryKey(obj, s.keyName)
}
