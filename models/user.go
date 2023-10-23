package models

// struct user represents data of user.
type User struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Mobile      string  `json:"mobile" binding:"required,e164"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"` // lat and long separate colums
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// foreign key of products id
