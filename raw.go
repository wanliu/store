package store

import (
	"github.com/boltdb/bolt"
)

func (s *store) Read(tx *bolt.Tx, id uint64) (obj interface{}, err error) {
	obj = s.NewEntity()

	if id == 0 {
		return nil, ErrZeroID
	}

	var (
		b     *bolt.Bucket
		buf   []byte
		pkey  Keyer
		bName = []byte(s.Name())
	)

	if b = tx.Bucket(bName); b == nil {
		return nil, ErrMissingBucket
	}

	if pkey, err = s.primaryKey(obj); err != nil {
		return nil, err
	}

	pkey.Set(id)

	buf = b.Get(s.IdName(pkey))

	if err = s.Decode(buf, obj); err != nil {
		return nil, err
	}

	return obj, err
}

func (s *store) Write(tx *bolt.Tx, obj interface{}) (err error) {
	var ok bool

	if ok, err = s.Validate(obj); !ok {
		return err
	}

	// write to boltdb buckets
	// return s.db.Update(func(tx *bolt.Tx) error {
	var (
		b     *bolt.Bucket
		pkey  Keyer
		bName = []byte(s.Name())
	)
	// Open or Create Bucket
	if b, err = tx.CreateBucketIfNotExists(bName); err != nil {
		return err
	}

	// Set Primary Key and return ID
	if pkey, err = s.primaryKey(obj); err != nil {
		return err
	}

	if pkey.Value() == 0 {
		if s.autoIncrement {
			id, _ := b.NextSequence()
			pkey.Set(id)
		} else {
			return ErrInvalidPKey
		}
	}

	// Call beforeCreate Hook

	// Timestamps
	createTimestamp(obj)

	// bckIndex, err := b.CreateBucketIfNotExists(s.indexName)
	// if err != nil {
	// 	return err
	// }

	// Encoder
	buf, err := s.Encode(obj)

	if err != nil {
		return err
	}

	if err = b.Put(s.IdName(pkey), buf); err != nil {
		return err
	}

	// write Indexes
	// if err = s.WriteIndexes(bckIndex, obj); err != nil {
	// 	return err
	// }

	// Call afterCreate Hook
	return nil
	// })
}

func (s *store) Delete(tx *bolt.Tx, id uint64) (err error) {
	var (
		b     *bolt.Bucket
		pkey  Keyer
		bName = []byte(s.Name())
	)

	obj, err := s.Read(tx, id)

	if err != nil {
		return err
	}
	// Open or Create Bucket
	if b = tx.Bucket(bName); b == nil {
		return ErrMissingBucket
	}

	// Set Primary Key and return ID
	if pkey, err = s.primaryKey(obj); err != nil {
		return err
	}

	if err = b.Delete(s.IdName(pkey)); err != nil {
		return err
	}

	return nil
}
