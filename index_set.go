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
	isNew   bool
	indexes []Index
}

func (s *store) NewIdxSet(obj interface{}, tx *bolt.Tx, b *bolt.Bucket) *indexSet {
	idxSet := &indexSet{obj: obj, s: s, tx: tx, b: b, isNew: true}
	return idxSet
}

func (s *store) LoadIdxSet(obj interface{}, tx *bolt.Tx, b *bolt.Bucket) *indexSet {
	idxSet := &indexSet{obj: obj, s: s, tx: tx, b: b, isNew: false}
	return idxSet
}

func (idxSet *indexSet) pure() {
	idxSet.indexes = nil
}

func (idxSet *indexSet) build() error {
	s := structs.New(idxSet.obj)
	pkey, err := idxSet.s.primaryKey(idxSet.obj)
	if err != nil {
		return err
	}

	for _, rule := range idxSet.s.indexes {
		feld, ok := s.FieldOk(rule.name)
		if !ok {
			return ErrorInvalidField{idxSet.obj, rule.name}
		}

		buf, err := idxSet.s.Encode(feld.Value())

		if err != nil {
			return err
		}

		idx := NewIndex(idxSet, idxSet.b, rule, buf, pkey.Value())
		idxSet.indexes = append(idxSet.indexes, idx)
	}

	return nil
}

func (idxSet *indexSet) reset() error {
	idxSet.pure()
	return idxSet.build()
}

func (idxSet *indexSet) Read() error {
	if err := idxSet.reset(); err != nil {
		return err
	}

	for _, idx := range idxSet.indexes {
		if err := idx.Load(); err != nil {
			return err
		}
	}

	return nil
}

func (idxSet *indexSet) Write() error {
	// if err := idxSet.reset(); err != nil {
	// 	return err
	// }

	pkey, err := idxSet.s.primaryKey(idxSet.obj)
	if err != nil {
		return err
	}

	for _, idx := range idxSet.indexes {
		if idxSet.isNew {
			idx.AddId(pkey.Value())
		}
		if err := idx.Save(); err != nil {
			return err
		}
	}

	return nil
}

func (idxSet *indexSet) Update(obj interface{}) error {
	s := structs.New(obj)

	for _, idx := range idxSet.indexes {
		feld, ok := s.FieldOk(idx.Name())
		if !ok {
			return ErrorInvalidField{obj, idx.Name()}
		}

		buf, err := idxSet.s.Encode(feld.Value())
		if err != nil {
			return err
		}

		idx.Update(buf)
	}

	return nil
}

func (idxSet *indexSet) Delete() error {
	for _, idx := range idxSet.indexes {
		if err := idx.Delete(); err != nil {
			return err
		}
	}

	return nil
}

func (idxSet *indexSet) Evaluate() error {
	for _, idx := range idxSet.indexes {
		if !idx.Evaluate() {
			return ErrUnqiueAwayls
		}
	}

	return nil
}
