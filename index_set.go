package store

import (
	"github.com/boltdb/bolt"
	"github.com/fatih/structs"
)

type indexSet struct {
	tx      *bolt.Tx
	s       *store
	b       *bolt.Bucket
	obj     interface{}
	indexes []Index
}

func (s *store) NewIdxSet(obj interface{}, tx *bolt.Tx, b *bolt.Bucket) *indexSet {
	idxSet := &indexSet{obj: obj, s: s, tx: tx, b: b}
	idxSet.build()
	return idxSet
}

func (idxSet *indexSet) build() error {
	s := structs.New(idxSet.obj)

	for _, rule := range idxSet.s.indexes {
		feld, ok := s.FieldOk(rule.name)
		if !ok {
			return ErrorInvalidField{idxSet.obj, rule.name}
		}

		buf, err := idxSet.s.Encode(feld.Value())

		if err != nil {
			return err
		}

		idx := NewIndex(idxSet, idxSet.b, rule, buf)
		idxSet.indexes = append(idxSet.indexes, idx)
	}

	return nil
}

func (idxSet *indexSet) Read() error {
	for _, idx := range idxSet.indexes {
		if err := idx.Load(); err != nil {
			return err
		}
	}

	return nil
}

func (idxSet *indexSet) Write() error {
	s := structs.New(idxSet.obj)

	for _, idx := range idxSet.indexes {
		feld, ok := s.FieldOk(idx.Name())
		if !ok {
			return ErrorInvalidField{idxSet.obj, idx.Name()}
		}

		buf, err := idxSet.s.Encode(feld.Value())
		if err != nil {
			return err
		}

		idx.Update(buf)
	}

	return nil
}

// func (idxSet *indexSet) Delete(obj interface{}) error {
// 	s := structs.New(obj)

// 	for _, idx := range idxSet.indexes {
// 		if err = idx.Delete(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
