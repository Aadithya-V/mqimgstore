package product

type AddableProduct struct {
	UserId             int      `json:"user_id" binding:"required"`
	ProductName        string   `json:"product_name" binding:"required"`
	ProductDescription string   `json:"product_description" binding:"required"` // text sql larger sizze
	ProductImages      []string `json:"product_images" binding:"required"`
	ProductPrice       float32  `json:"product_price" binding:"required"`
}
