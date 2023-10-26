package users

import (
	"database/sql"
)

type Respository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Respository {
	return &Respository{db}
}

func (r *Respository) GetUserByID(id int64) (User, error) {
	return userByID(id, r.db)
}
