package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IService interface {
	AddProduct(addableProduct AddableProduct) error
}

type Controller struct {
	s IService
}

func NewController(s IService) *Controller {
	return &Controller{s}
}

func (c *Controller) AddProduct(ctx *gin.Context) {
	var addableProduct AddableProduct
	if err := ctx.BindJSON(&addableProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, &gin.H{"error": "product information not submitted correctly. FORMAT: ..." + err.Error()}) // print json format for reference
		return
	}

	err := c.s.AddProduct(addableProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &gin.H{"error": "Unable to add product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, &gin.H{"message": "Product successfully added to database. Tip: It takes 10-20 minutes for uploaded images to become visible uniformly around the world. (Distribution to CDN)"})
}
