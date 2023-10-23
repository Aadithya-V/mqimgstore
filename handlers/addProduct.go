package handlers

import (
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aadithya-V/mqimgstore/models"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AddProduct(db *sql.DB, mqConn *amqp.Connection) func(ctx *gin.Context) {

	mqCh, err := mqConn.Channel() // TODO: check if the deferred mqConn.close at main also closes mqCh at the end of main to be safer.
	FailOnError(err, "Failed to open a channel")

	q, err := mqCh.QueueDeclare(
		"imgcompressionservice", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	// CLOSURE TO RETURN TO GIN
	fx := func(ctx *gin.Context) {
		var addProduct models.AddableProduct
		if err := ctx.BindJSON(&addProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, &gin.H{"error": "product information not submitted correctly. FORMAT: ..." + err.Error()}) // print json format for reference
			return
		}

		// Check if user_id is valid / already exists
		_, err := userByID(addProduct.UserId, db)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, &gin.H{"error": err.Error()})
			return
		}

		// FUTURE: Check if product exists and updation. OR add new and delete old prod.

		// add details to Products table
		productId, err := insertProduct(addProduct, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "Unable to add product: " + err.Error()})
			return
		}
		// newtask.go
		// Producer. Enqueue to RabbitMQ
		mqctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := make([]byte, 8)
		binary.LittleEndian.PutUint64(body, uint64(productId))

		err = mqCh.PublishWithContext(mqctx,
			"",     // exchange
			q.Name, // routing key
			true,   // mandatory // TODO: add NotifyReturn
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})
		if err != nil { // need to handle err better
			ctx.JSON(http.StatusInternalServerError, &gin.H{"error": "Added product but unable to process product images: " + err.Error()})
			return
		}

		// Retry failed enqueues in an exponentially increasing time with limit. Notify admin. Queue in a separate go queue / write to disk.
		// flush buffer by closing
		// return notify listener

		log.Printf(" [x] Sent %s\n", body)

		ctx.JSON(http.StatusCreated, &gin.H{"message": "Product successfully added to database. Tip: It takes 10-20 minutes for uploaded images to become visible uniformly around the world. (Distribution to CDN)"})
	}
	return fx
}

// move to sql helpers

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

// insertProduct adds the specified product to the database,
// returning the product_id of the new entry
func insertProduct(addProduct models.AddableProduct, db *sql.DB) (int64, error) {
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
