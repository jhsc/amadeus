// Package store provides an amadeus data store interface.
package store

import (
	"errors"
)

var (
	// ErrNotFound means the requested item is not found.
	ErrNotFound = errors.New("store: item not found")
	// ErrConflict means the operation failed because of a conflict between items.
	ErrConflict = errors.New("store: item conflict")
)

// Store is an amadadeus data store interface.
type Store interface {
	Projects() ProjectStore
}

// ProjectStore is an amadeus project data store interface.
type ProjectStore interface {
	New(projectName string) (int64, error)
	Get(id int64) (*Project, error)
	GetMany(ids []int64) (map[int64]*Project, error)
}
