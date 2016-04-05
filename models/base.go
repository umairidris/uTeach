// Package models contains models to abstract the db.
// All models methods have a tx as the first parameter. This allows handlers that use multiple calls / tables to
// use a single transaction and then commit that so that all actions are committed or none are.
// If a handler is using a single action, they can simply pass in nil for the tx.
package models

import (
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
)

type Base struct {
	db *sqlx.DB
}

func (b *Base) Exec(tx *sqlx.Tx, query string, args ...interface{}) (driver.Result, error) {

	if tx != nil {
		return tx.Exec(query, args...)
	}

	return b.db.Exec(query, args...)
}

func (b *Base) Get(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Get(dest, query, args...)
	}

	return b.db.Get(dest, query, args...)
}

func (b *Base) Select(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Select(dest, query, args...)
	}

	return b.db.Select(dest, query, args...)
}

func (b *Base) Query(tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Rows, error) {
	if tx != nil {
		return tx.Queryx(query, args...)
	}

	return b.db.Queryx(query, args...)
}
