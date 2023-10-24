package user

import (
	"database/sql"

	"github.com/Aadithya-V/mqimgstore/models"
)

type Respository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Respository {
	return &Respository{db}
}

func (r *Respository) GetUserByID(id int) (models.User, error) {
	return userByID(id, r.db)
}
