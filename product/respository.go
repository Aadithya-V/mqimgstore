package product

import (
	"database/sql"
)

type Respository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Respository {
	return &Respository{db}
}

func (r *Respository) InsertProduct(addProduct AddableProduct) (int64, error) {
	return insertProduct(addProduct, r.db)
}
