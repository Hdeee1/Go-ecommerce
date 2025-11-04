package middleware

import (
	"net/http"
	"strings"

	"github.com/Hdeee1/go-ecommerce/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.GetHeader("Authorization")

		if clientToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			ctx.Abort()
			return 
		}

		clientToken = strings.TrimPrefix(clientToken, "Bearer ")

		claims, err := tokens.ValidateToken(clientToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return 
		}

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.First_Name)
		ctx.Set("last_name", claims.Last_Name)
		ctx.Set("user_id", claims.User_ID)

		ctx.Next()
	}
}