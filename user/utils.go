package user

import (
	"database/sql"
	"fmt"

	"github.com/Aadithya-V/mqimgstore/models"
)

// userByID queries for the user with the specified ID.
func userByID(id int, db *sql.DB) (models.User, error) {
	// A models.User to hold data from the returned row.
	var user models.User

	row := db.QueryRow("SELECT * FROM users WHERE user_id = ? ", id)
	if err := row.Scan(&user.ID, &user.Name, &user.Mobile, &user.Latitude, &user.Longitude, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("userById %d: no such user", id)
		}
		return user, fmt.Errorf("userById %d: %v", id, err)
	}
	return user, nil
}
