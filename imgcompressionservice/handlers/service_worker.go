package handlers

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ServiceWorker(d amqp.Delivery, db *sql.DB) {

	productID := int64(binary.LittleEndian.Uint64(d.Body))
	fmt.Println(productID)

	uncompUrls, err := retrieveUNcompressedImgURLs(productID, db, d)
	if err != nil {
		fmt.Printf("Error during product_images retrieval: %v", err.Error())
		return
	}

	fmt.Println(uncompUrls)

	d.Ack(false)

}

// d ampq.Delivery is a paramenter for Negative Ack on failure: either drop or retry message
func retrieveUNcompressedImgURLs(productID int64, db *sql.DB, d amqp.Delivery) ([]string, error) {
	var urls []byte
	var res []string

	row := db.QueryRow("SELECT product_images FROM products WHERE product_id = ? ", productID)
	if err := row.Scan(&urls); err != nil {
		if err == sql.ErrNoRows {
			d.Nack(false, false) // drop message from queue
			return res, fmt.Errorf("product_id %d: no such product", productID)
		}
		d.Nack(false, true) // retry
		return res, fmt.Errorf("product_id %d: %v", productID, err)
	}
	json.Unmarshal(urls, &res)
	return res, nil
}

/* func updateCompressedImgURLs(productID int, jsonObj []byte, db *sql.DB) error {

}

func CompressWorker() {

}

func storeImg(fd) {

}
*/
