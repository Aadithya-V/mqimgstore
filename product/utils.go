package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// insertProduct adds the specified product to the database,
// returning the product_id of the new entry
func insertProduct(addProduct AddableProduct, db *sql.DB) (int64, error) {
	jsonObjImgLinks, err := json.Marshal(addProduct.ProductImages)
	if err != nil {
		return 0, fmt.Errorf("insertProduct: %v", err)
	}

	result, err := db.Exec("INSERT INTO products (user_id, product_name, product_description, product_images, product_price) VALUES (?, ?, ?, ?, ?)",
		addProduct.UserId, addProduct.ProductName, addProduct.ProductDescription, jsonObjImgLinks, addProduct.ProductPrice)
	if err != nil {
		return 0, fmt.Errorf("insertProduct: %v", err)
	}
	ProductId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertProduct: %v", err)
	}
	return ProductId, nil
}
