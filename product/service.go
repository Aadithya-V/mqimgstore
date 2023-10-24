package product

import (
	"context"
	"encoding/binary"
	"log"
	"time"

	"github.com/Aadithya-V/mqimgstore/models"
	"github.com/Aadithya-V/mqimgstore/queue"
)

type IProductRepository interface {
	InsertProduct(addProduct AddableProduct) (int64, error)
}

type IUserRepository interface {
	GetUserByID(id int) (models.User, error)
}

type Service struct {
	pRepo IProductRepository
	uRepo IUserRepository
	msgB  *queue.MessageBroker
}

func NewService(pRepo IProductRepository, uRepo IUserRepository, msgB *queue.MessageBroker) *Service {
	return &Service{pRepo, uRepo, msgB}
}

func (s *Service) AddProduct(addableProduct AddableProduct) error {
	// Check if user_id is valid / already exists
	_, err := s.uRepo.GetUserByID(addableProduct.UserId)
	if err != nil {
		return err
	}

	// Add new product to db
	prodID, err := s.pRepo.InsertProduct(addableProduct)
	if err != nil {
		return err
	}

	// Publish to MessageBroker's imgcompressionservice Queue <- prodID
	body := make([]byte, 8)
	binary.LittleEndian.PutUint64(body, uint64(prodID))

	mqctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.msgB.Publish(mqctx, queue.ImageCompressionQueue, body)
	if err != nil {
		// Retry failed enqueues in an exponentially increasing time with limit. Notify admin. Queue in a separate go queue / write to disk.
		// flush buffer by closing
		// return notify listener

		return err
	}

	log.Printf(" [x] Sent %s\n", body)

	return nil
}
