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

	app := controllers.NewApplication(
		userRepo,
		productRepo,
		orderRepo,
	)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router, app)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.POST("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}