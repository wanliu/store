package store

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/fatih/structs"
)

const PriKeyName = "ID"

var bSep = []byte(":")

type indexRule struct {
	s       *store
	keyName string
	name    string
	idxType string
}

type index struct {
	rule     indexRule
	set      *indexSet
	b        *bolt.Bucket
	dirty    bool
	value    []byte
	elements []uint64
}

type Index interface {
	Evaluate() bool
	Elements() []uint64
	// Classic() []byte
	// Read(idx []byte) []byte
	Load() error
	Update(idx []byte) error
	Delete() error
	Name() string
}

// func RegIndex(name string, idx Index) {

// }

type uniqueIndex struct {
	index
}

func NewIndex(idxSet *indexSet, b *bolt.Bucket, rule indexRule, value []byte) Index {
	switch rule.idxType {
	case "index":
		return &index{set: idxSet, b: b, rule: rule, value: value}
	case "unique":
		return &uniqueIndex{index{set: idxSet, b: b, rule: rule, value: value}}
	}
	return nil
}

// index
func (idx *index) idxName() ([]byte, error) {
	bName := []byte(idx.rule.name)

	return bytes.Join([][]byte{bName, idx.value}, bSep), nil
}

func (idx *index) Evaluate() bool {
	return true
}

func (idx *index) Name() string {
	return idx.rule.name
}

func (idx *index) Load() error {
	return idx.read()
}

func (idx *index) Update(newVal []byte) (err error) {
	if bytes.Compare(idx.value, newVal) == 0 {
		return nil
	}

	idx.dirty = true

	idx.removeId(idx.targetId())
	if err = idx.write(); err != nil {
		return err
	}

	newIdx := index{b: idx.b, rule: idx.rule, value: newVal}
	newIdx.Load()
	newIdx.addId(idx.targetId())

	if err = newIdx.write(); err != nil {
		return err
	}
	return nil
}

func (idx *index) Delete() (err error) {
	if err = idx.read(); err != nil {
		return err
	}

	idx.removeId(idx.targetId())
	if len(idx.elements) > 0 {
		return ErrNotEmptyIndex
	} else {
		return idx.pure()
	}
}

func (idx *index) Elements() []uint64 {
	return idx.elements
}

// func (idx *index) String() string {
// 	return idx.rule.name
// }

func (idx *index) read() error {
	var name, _ = idx.idxName()
	v := idx.b.Get(name)
	err := idx.rule.s.Decode(v, &idx.elements)
	if err != nil {
		return err
	}

	idx.dirty = false
	return nil
}

func (idx *index) write() error {
	if len(idx.elements) > 0 {
		return idx._write()
	} else {
		return idx.pure()
	}
}

func (idx *index) _write() error {
	v, err := idx.rule.s.Encode(idx.elements)
	if err != nil {
		return err
	}
	var name, _ = idx.idxName()

	if err = idx.b.Put(name, v); err != nil {
		return err
	}
	idx.dirty = false
	return nil
}

func (idx *index) addId(id uint64) int {
	var (
		idxId uint64
	)

	for _, idxId = range idx.elements {
		if idxId == id {
			return 0
		}
	}

	idx.elements = append(idx.elements, id)
	return 1
}

func (idx *index) removeId(id uint64) int {
	for i, idxId := range idx.elements {
		if idxId == id {
			idx.elements = append(idx.elements[:i], idx.elements[i+1:]...)
			return 1
		}
	}

	return 0
}

func (idx *index) targetId() uint64 {
	pkey, err := idx.set.s.primaryKey(idx.set.obj)
	if err != nil {
		return 0
	}

	return uint64(pkey.Value())
}

func (idx *index) pure() error {
	if len(idx.elements) == 0 {
		if err := idx.b.Delete(idx.value); err != nil {
			return err
		}
		idx.dirty = false
		return nil
	} else {
		return ErrNotEmptyIndex
	}
}

// uniqueIndex
func (idx *uniqueIndex) Evaluate() bool {
	name, _ := idx.idxName()

	if len(name) == 0 {
		return false
	}

	if err := idx.read(); err != nil {
		return false
	}

	for _, id := range idx.elements {
		if idx.targetId() == id {
			return true
		}
	}
	return false
}

// func (idx *uniqueIndex) Write(b *bolt.Bucket, obj interface{}) error {
// 	if idx.Evaluate(b, obj) {
// 		return idx.index.Write(b, obj)
// 	} else {
// 		return ErrUnqiueAwayls
// 	}
// }

// store index methods
func (s *store) BuildIndexes(obj interface{}) (indexes []indexRule, err error) {
	st := structs.New(obj)

	for _, feld := range st.Fields() {
		idxType := feld.Tag("index")

		switch idxType {
		case "index":
			idx := indexRule{s: s, name: feld.Name(), keyName: PriKeyName, idxType: "index"}
			indexes = append(indexes, idx)
		case "unique":
			idx := indexRule{s: s, name: feld.Name(), keyName: PriKeyName, idxType: "unique"}
			indexes = append(indexes, idx)
		default:
			continue
		}
	}

	return indexes, nil
}

// func (s *store) EvaluateIndexes(b *bolt.Bucket, obj interface{}) error {
// 	for _, idx := range s.elements {
// 		if !idx.Evaluate(b, obj) {
// 			name := idx.Name()
// 			st := structs.New(obj)
// 			if feld, ok := st.FieldOk(name); ok {
// 				return &ErrorUniqueIndex{s, idx, feld.Value()}
// 			} else {
// 				return &ErrorInvalidField{obj, name}
// 			}
// 		}
// 	}

// 	return nil
// }

// func (s *store) WriteIndexes(b *bolt.Bucket, obj interface{}) (err error) {
// 	if err = s.EvaluateIndexes(b, obj); err != nil {
// 		return err
// 	}

// 	for _, idx := range s.elements {
// 		if err = idx.Write(b, obj); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *store) RemoveIndexes(b *bolt.Bucket, obj interface{}) (err error) {
// 	for _, idx := range s.elements {
// 		if err = idx.Delete(b, obj); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
