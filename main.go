package main

import (
	"log"

	"github.com/Hdeee1/go-ecommerce/config"
	"github.com/Hdeee1/go-ecommerce/controllers"
	"github.com/Hdeee1/go-ecommerce/database"
	"github.com/Hdeee1/go-ecommerce/middleware"
	"github.com/Hdeee1/go-ecommerce/repository"
	"github.com/Hdeee1/go-ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	port := cfg.Port

	database.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)
	addressRepo := repository.NewAddressRepository(database.DB)

	app := controllers.NewApplication(
		userRepo,
		productRepo,
		orderRepo,
		addressRepo,
	)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router, app)
	router.Use(middleware.Authentication())

	routes.ProductRoutes(router, app)
	log.Fatal(router.Run(":" + port))
}