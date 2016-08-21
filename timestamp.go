package store

import (
	"time"

	"github.com/fatih/structs"
)

type Timestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func createTimestamp(obj interface{}) error {
	s := structs.New(obj)
	if created_at, ok := s.FieldOk("CreatedAt"); ok {
		created_at.Set(time.Now())
	}

	if updated_at, ok := s.FieldOk("UpdatedAt"); ok {
		updated_at.Set(time.Now())
	}

	return nil
}

func updateTimestamp(obj interface{}) error {
	s := structs.New(obj)

	if updated_at, ok := s.FieldOk("UpdatedAt"); ok {
		updated_at.Set(time.Now())
	}

	return nil
}
