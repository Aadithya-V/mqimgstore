package server

import (
	"github.com/Aadithya-V/mqimgstore/product"
	"github.com/Aadithya-V/mqimgstore/user"
)

func (s *Server) Product() {
	productRepository := product.NewRepository(s.db)
	userRepository := user.NewRepository(s.db)

	productService := product.NewService(productRepository, userRepository, s.msgB)

	productController := product.NewController(productService)

	s.router.POST("/product", productController.AddProduct)

}
