package product_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aadithya-V/mqimgstore/mocks"
	"github.com/Aadithya-V/mqimgstore/product"
	"github.com/Aadithya-V/mqimgstore/queue"
	"github.com/Aadithya-V/mqimgstore/users"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* type UnitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
} */

func TestAddProduct(t *testing.T) {

	prod := product.AddableProduct{
		UserId:             2,
		ProductName:        "XYZ",
		ProductDescription: "Loren Ipsum",
		ProductImages:      []string{"url", "url/2"},
		ProductPrice:       10.00,
	}

	mockpRepo := mocks.NewIProductRepository(t)
	mockpRepo.On("InsertProduct", prod).Return(int64(1), nil)

	mockuRepo := mocks.NewIUserRepository(t)
	mockuRepo.On("GetUserByID", prod.UserId).Return(users.User{
		ID:          2,
		Name:        "John Doe",
		Description: "test user adding product named XYZ",
	}, nil)

	mockMsgB := mocks.NewIMsgB(t)
	mockMsgB.On("Publish", mock.AnythingOfType("*context.timerCtx"), queue.ImageCompressionQueue, mock.AnythingOfType("[]uint8")).Return(nil) //  TODO: bytes/encoding size []int8 increases as sizd of productid increases.

	mockService := product.NewService(mockpRepo, mockuRepo, mockMsgB)

	productController := product.NewController(mockService)

	router := gin.Default()
	router.POST("/product", productController.AddProduct)

	body, _ := json.Marshal(prod)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code, w.Body)

	body = []byte("") // TODO: Test Passing emmpty / body json with not all required product fields
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/product", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code, w.Body)

}
