package store

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/fatih/structs"
)

const PriKeyName = "ID"

var bSep = []byte(":")

type index struct {
	s       *store
	keyName string
	name    string
}

type Index interface {
	Evaluate(*bolt.Bucket, interface{}) bool
	Classic() []byte
	Read(*bolt.Bucket, interface{}) []byte
	Write(*bolt.Bucket, interface{}) error
	Delete(*bolt.Bucket, interface{}) error
	Name() string
}

// func RegIndex(name string, idx Index) {

// }

type uniqueIndex struct {
	index
}

// index
func (idx *index) idxName(obj interface{}) ([]byte, error) {
	s := structs.New(obj)
	feld := s.Field(idx.name)
	bName := []byte(idx.name)
	bVal, err := idx.s.Encode(feld.Value())

	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{bName, bVal}, bSep), nil
}

func (idx *index) Evaluate(b *bolt.Bucket, obj interface{}) bool {
	return true
}

func (idx *index) Name() string {
	return idx.name
}

func (idx *index) Read(b *bolt.Bucket, obj interface{}) []byte {
	name, _ := idx.idxName(obj)
	return b.Get(name)
}

func (idx *index) Write(b *bolt.Bucket, obj interface{}) error {
	name, _ := idx.idxName(obj)
	pkey, _ := PrimaryKey(obj, idx.keyName)
	idxValue := pkey.Value()
	bVal, err := idx.s.Encode(idxValue)

	if err != nil {
		return err
	}

	return b.Put(name, bVal)
}

func (idx *index) Delete(b *bolt.Bucket, obj interface{}) error {
	name, _ := idx.idxName(obj)

	return b.Delete(name)
}

func (idx *index) String() string {
	return idx.name
}

// uniqueIndex
func (idx *uniqueIndex) Evaluate(b *bolt.Bucket, obj interface{}) bool {
	name, _ := idx.idxName(obj)
	pkey, err := PrimaryKey(obj, idx.keyName)
	if err != nil {
		return false
	}

	if len(name) == 0 {
		return false
	}

	v := b.Get(name)
	if len(v) > 0 {
		var id uint64
		if err = idx.s.Decode(v, &id); err != nil {
			return false
		}

		if id == pkey.Value() {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

// func (idx *uniqueIndex) Write(b *bolt.Bucket, obj interface{}) error {
// 	if idx.Evaluate(b, obj) {
// 		return idx.index.Write(b, obj)
// 	} else {
// 		return ErrUnqiueAwayls
// 	}
// }

// store index methods
func (s *store) BuildIndexes(obj interface{}) (indexes []Index, err error) {
	st := structs.New(obj)

	for _, feld := range st.Fields() {
		idxType := feld.Tag("index")

		switch idxType {
		case "index":
			idx := index{s: s, name: feld.Name(), keyName: PriKeyName}
			indexes = append(indexes, &idx)
		case "unique":
			idx := uniqueIndex{index{s: s, name: feld.Name(), keyName: PriKeyName}}
			indexes = append(indexes, &idx)
		default:
			continue
		}
	}

	return indexes, nil
}

func (s *store) EvaluateIndexes(b *bolt.Bucket, obj interface{}) error {
	for _, idx := range s.indexes {
		if !idx.Evaluate(b, obj) {
			name := idx.Name()
			st := structs.New(obj)
			if feld, ok := st.FieldOk(name); ok {
				return &ErrorUniqueIndex{s, idx, feld.Value()}
			} else {
				return &ErrorInvalidField{obj, name}
			}
		}
	}

	return nil
}

func (s *store) WriteIndexes(b *bolt.Bucket, obj interface{}) (err error) {
	if err = s.EvaluateIndexes(b, obj); err != nil {
		return err
	}

	for _, idx := range s.indexes {
		if err = idx.Write(b, obj); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) RemoveIndexes(b *bolt.Bucket, obj interface{}) (err error) {
	for _, idx := range s.indexes {
		if err = idx.Delete(b, obj); err != nil {
			return err
		}
	}

	return nil
}
