package models

type Product struct {
	ProductId               int      `json:"product_id"`
	ProductName             string   `json:"product_name" binding:"required"`
	ProductDescription      string   `json:"product_description" binding:"required"` // text sql larger sizze
	ProductImages           []string `json:"product_images" binding:"required"`
	ProductPrice            float32  `json:"product_price" binding:"required"`
	CompressedProductImages []string
	CreatedAt               string
	UpdatedAt               string
}

type AddableProduct struct {
	UserId             int      `json:"user_id" binding:"required"`
	ProductName        string   `json:"product_name" binding:"required"`
	ProductDescription string   `json:"product_description" binding:"required"` // text sql larger sizze
	ProductImages      []string `json:"product_images" binding:"required"`
	ProductPrice       float32  `json:"product_price" binding:"required"`
}

// url prod imgs regex validation

// foreign key of users id
