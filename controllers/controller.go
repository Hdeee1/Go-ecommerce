package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/gin-gonic/gin"
)

func HashPassword(password string) string {
	
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return 
		}

		UserCollection.CountDocument(ctx, )
	}
}

func Login() gin.HandlerFunc {
	
}

func ProductViewerAdmin() gin.HandlerFunc {
	
}

func SearchProduct() gin.HandlerFunc {
	
}

func SearchProductByQuery() gin.HandlerFunc {
	
}