package server

import (
	"github.com/Aadithya-V/mqimgstore/product"
	"github.com/Aadithya-V/mqimgstore/users"
)

func (s *Server) Product() {
	productRepository := product.NewRepository(s.db)
	userRepository := users.NewRepository(s.db)

	productService := product.NewService(productRepository, userRepository, s.msgB)

	productController := product.NewController(productService)

	s.router.POST("/product", productController.AddProduct)

}
