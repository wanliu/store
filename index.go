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
	value    []byte
	target   uint64
	elements []uint64
}

type Index interface {
	Evaluate() bool
	Elements() []uint64
	// Classic() []byte
	// Read(idx []byte) []byte
	AddId(uint64) int
	RemoveId(uint64) int
	Load() error
	Update(newIdx []byte) error
	Save() error
	Delete() error
	Name() string
}

// func RegIndex(name string, idx Index) {

// }

type uniqueIndex struct {
	index
}

func NewIndex(idxSet *indexSet, b *bolt.Bucket, rule indexRule, value []byte, target uint64) Index {
	switch rule.idxType {
	case "index":
		return &index{set: idxSet, b: b, rule: rule, value: value, target: target}
	case "unique":
		return &uniqueIndex{index{set: idxSet, b: b, rule: rule, value: value, target: target}}
	}
	return nil
}

// index
func (idx *index) idxName() ([]byte, error) {
	if idx.value == nil {
		return nil, ErrMissingValue
	}

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

	// idx.dirty = true

	idx.RemoveId(idx.target)
	if err = idx.Save(); err != nil {
		return err
	}

	newIdx := NewIndex(idx.set, idx.b, idx.rule, newVal, idx.target)
	newIdx.Load()
	newIdx.AddId(idx.target)

	if err = newIdx.Save(); err != nil {
		return err
	}
	return nil
}

func (idx *index) Save() (err error) {
	return idx.write()
}

func (idx *index) Delete() (err error) {
	if err = idx.read(); err != nil {
		return err
	}

	idx.RemoveId(idx.target)
	if len(idx.elements) > 0 {
		return idx._write()
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

func (idx *index) read() (err error) {
	var name, _ = idx.idxName()
	v := idx.b.Get(name)

	if len(v) > 0 {
		err = idx.rule.s.Decode(v, &idx.elements)
	} else {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

//  write or clean indexes
func (idx *index) write() error {
	if len(idx.elements) > 0 {
		return idx._write()
	} else {
		return idx.pure()
	}
}

// _write indexes
func (idx *index) _write() error {
	v, err := idx.rule.s.Encode(idx.elements)
	if err != nil {
		return err
	}
	var name, _ = idx.idxName()

	if err = idx.b.Put(name, v); err != nil {
		return err
	}

	return nil
}

func (idx *index) AddId(id uint64) int {
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

func (idx *index) RemoveId(id uint64) int {
	for i, idxId := range idx.elements {
		if idxId == id {
			idx.elements = append(idx.elements[:i], idx.elements[i+1:]...)
			return 1
		}
	}

	return 0
}

func (idx *index) pure() error {
	if len(idx.elements) == 0 {
		name, _ := idx.idxName()

		if err := idx.b.Delete(name); err != nil {
			return err
		}
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

	if len(idx.elements) == 0 {
		return true
	}

	for _, id := range idx.elements {
		if idx.target == id {
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
