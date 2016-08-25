package store

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fatih/structs"
	"github.com/imdario/mergo"
	"github.com/ugorji/go/codec"
)

var db *bolt.DB

func Open(dbpath string) (*bolt.DB, error) {
	var err error
	db, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	return db, err
}

func NewStore(entity interface{}, opts Map) Store {
	s := &store{Entity: entity}

	if opts == nil {
		opts = defaultOption()
	} else {
		mergo.Merge(&opts, defaultOption())
	}

	if opts["codec"] != nil {
		s.codec = opts["codec"].(codec.Handle)
	}

	if opts["db"] != nil {
		s.db = opts["db"].(*bolt.DB)
	}

	if opts["key"] != nil {
		s.keyName = opts["key"].(string)
	}

	if opts["auto_increment"] != nil {
		s.autoIncrement = opts["auto_increment"].(bool)
	}

	if opts["index_name"] != nil {
		s.indexName = opts["index_name"].([]byte)
	}

	s.indexes, _ = s.BuildIndexes(entity)

	return s
}

func defaultOption() Map {
	return Map{
		"key":            "ID",
		"codec":          defCodec,
		"db":             db,
		"auto_increment": true,
		"index_name":     []byte("indexes"),
	}
}

func (s *store) NewEntity() interface{} {
	v := reflect.ValueOf(s.Entity)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	o := reflect.New(t)

	return o.Interface()
}

func (s *store) Create(attr Map) (obj interface{}, err error) {
	// instance objecct
	obj = s.NewEntity()
	err = s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.Name()))
		if err != nil {
			return err
		}

		// merge attributes
		if err = MergeStruct(obj, attr); err != nil {
			return err
		}

		bIdx, err := b.CreateBucketIfNotExists(IdxBucket)
		if err != nil {
			return err
		}

		idxSet := s.NewIdxSet(obj, tx, bIdx)
		idxSet.Read()

		defer func() {
			if err == nil {
				err = idxSet.Write()
			}
		}()

		if err = idxSet.Evaluate(); err != nil {
			return err
		}

		err = s.Write(tx, obj)
		return err
	})

	return obj, err
}

func (s *store) Get(id uint64) (obj interface{}, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		obj, err = s.Read(tx, id)
		return err
	})

	return obj, err
}

func (s *store) Put(id uint64, attr Map) (err error) {
	return s.db.Update(func(tx *bolt.Tx) error {
		obj, err := s.Read(tx, id)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists([]byte(s.Name()))
		if err != nil {
			return err
		}
		bIdx, err := b.CreateBucketIfNotExists(IdxBucket)
		if err != nil {
			return err
		}

		idxSet := s.NewIdxSet(obj, tx, bIdx)
		idxSet.Read()

		defer func() {
			if err == nil {
				err = idxSet.Update(obj)
			}
		}()

		if err = MergeStruct(obj, attr); err != nil {
			return err
		}

		evlSet := s.NewIdxSet(obj, tx, bIdx)
		evlSet.Read()
		if err = evlSet.Evaluate(); err != nil {
			return err
		}

		pkey, err := s.primaryKey(obj)

		if err != nil {
			return err
		}

		pkey.Set(id)

		err = s.Write(tx, obj)
		return err
	})
}

func (s *store) Remove(id uint64) (err error) {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.Name()))
		if err != nil {
			return err
		}

		bIdx, err := b.CreateBucketIfNotExists(IdxBucket)
		if err != nil {
			return err
		}

		obj, err := s.Read(tx, id)
		if err != nil {
			return err
		}

		idxSet := s.NewIdxSet(obj, tx, bIdx)
		idxSet.Read()

		defer func() {
			if err == nil {
				err = idxSet.Delete()
			}
		}()
		err = s.Delete(tx, id)
		return err
	})
}

func (s *store) Name() string {
	if len(s.bucketName) > 0 {
		return s.bucketName
	}
	s.bucketName = structs.Name(s.Entity)
	return s.bucketName
}

func (s *store) String() string {
	return s.Name() + "Store"
}

// func printElements(idxSet *indexSet) {
// 	for _, idx := range idxSet.indexes {
// 		log.Printf("elements  %+v", idx.Elements())
// 	}
// }
