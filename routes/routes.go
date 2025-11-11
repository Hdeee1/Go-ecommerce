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

func ProductRoutes(incomingRoutes *gin.Engine, app *controllers.Application) {
	authGroup := incomingRoutes.Group("/")
	{
		authGroup.POST("/admin/addproduct", controllers.ProductViewerAdmin(app))

		authGroup.GET("/addtocart", app.AddToCart())
		authGroup.GET("/removeitem", app.RemoveItem())
		authGroup.GET("/cartcheckout", app.BuyFromCart())
		authGroup.POST("/instantbuy", app.InstantBuy())

		addressGroup := authGroup.Group("/users/address")
		{
			addressGroup.POST("", app.AddAddress())
			addressGroup.PUT("/home", app.EditHomeAddress())
			addressGroup.PUT("/work", app.EditWorkAddress())
			addressGroup.DELETE("/:address_id", app.DeleteAddress())
		}
	}
}