package store

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/ugorji/go/codec"
)

var IdxBucket = []byte("Indexes")

type store struct {
	keyName       string
	bucketName    string
	db            *bolt.DB
	autoIncrement bool
	codec         codec.Handle
	indexes       []indexRule
	indexName     []byte
	// eacher    *Eacher
	Entity interface{}
	// Relations []Relation
}

type Map map[string]interface{}

type Store interface {
	Entitier
	Storage
}

type Entitier interface {
	NewEntity() interface{}
}

type Storage interface {
	Create(Map) (interface{}, error)
	Get(uint64) (interface{}, error)
	Put(uint64, Map) error
	Remove(uint64) error
	// Read(uint64) error
	// Write(interface{}) error
	// Delete(uint64) error
}

type AutoIncrement interface {
	NextId() uint64
}

type Encoder interface {
	Encode(buf *bytes.Buffer) error
}

type Decoder interface {
	Decode([]byte) error
}
