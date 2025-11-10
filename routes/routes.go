package routes

import (
	"github.com/Hdeee1/go-ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, app *controllers.Application) {
	incomingRoutes.POST("/users/signup", controllers.SignUp(app))
	incomingRoutes.POST("/users/login", controllers.Login(app))

	incomingRoutes.GET("/users/productview", controllers.SearchProduct(app))
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery(app))
}
