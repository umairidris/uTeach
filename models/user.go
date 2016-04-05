package models

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{Base{db}}
}

type User struct {
	ID      int64
	Email   string
	Name    string
	IsAdmin bool `db:"is_admin"`
}

type UserModel struct {
	Base
}

func (um *UserModel) GetUserByID(tx *sqlx.Tx, id int64) (*User, error) {
	user := &User{}
	err := um.Get(tx, user, "SELECT * FROM users WHERE id=?", id)
	return user, err
}

// GetUserByEmail returns record by email.
func (um *UserModel) GetUserByEmail(tx *sqlx.Tx, email string) (*User, error) {
	user := &User{}
	err := um.Get(tx, user, "SELECT * FROM users WHERE email=?", email)
	return user, err
}

// Signup creates a new record of user.
func (um *UserModel) Signup(tx *sqlx.Tx, email, name string) (*User, error) {
	if email == "" || name == "" {
		return nil, errors.New(".")
	}

	email = strings.ToLower(email)
	name = strings.Title(name)
	_, err := um.Exec(tx, "INSERT INTO users(email, name) VALUES(?, ?)", email, name)
	if err != nil {
		return nil, err
	}

	return um.GetUserByEmail(tx, email)
}
