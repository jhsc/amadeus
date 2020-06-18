// Package mysql provides a MySQL implementation of the amadeus data store interface.
package mysql

import (
	"bytes"
	"database/sql"
	"fmt"

	// "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"gitlab.com/jhsc/amadeus/store"
)

// Store is a mysql implementation of store.
type Store struct {
	db           *sql.DB
	projectStore *projectStore
}

// Projects returns a project store.
func (s *Store) Projects() store.ProjectStore {
	return s.projectStore
}

// Connect connects to a store.
func Connect() (*Store, error) {
	// MYSQL
	// connstr := fmt.Sprintf(
	// 	"%s:%s@tcp(%s)/%s?parseTime=true",
	// 	username, password, address, database,
	// )

	// db, err := sql.Open("mysql", connstr)
	sqlite3.Version()
	db, err := sql.Open("sqlite3", "./store.db")

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{
		db:           db,
		projectStore: &projectStore{db: db},
	}

	err = s.Migrate()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Migrate migrates the store database.
func (s *Store) Migrate() error {
	for _, q := range migrate {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

// Drop drops the store database.
func (s *Store) Drop() error {
	for _, q := range drop {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

// Reset resets the store database.
func (s *Store) Reset() error {
	err := s.Drop()
	if err != nil {
		return err
	}
	return s.Migrate()
}

type scanner interface {
	Scan(v ...interface{}) error
}

func placeholders(count int) string {
	buf := new(bytes.Buffer)
	for i := 0; i < count; i++ {
		buf.WriteByte('?')
		if i < count-1 {
			buf.WriteByte(',')
		}
	}
	return buf.String()
}

// func isUniqueConstraintError(err error) bool {
// 	if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1062 {
// 		return true
// 	}
// 	return false
// }
