package controllers

import (
	"net/http"

	"github.com/Hdeee1/go-ecommerce/database"
	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	
// }

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// check email duplication
		var existingUser models.User
		err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist"})
			return 
		}

		err = database.DB.Where("phone = ?", user.Phone).First(&existingUser).Error
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone already exist"})
			return
		}

		// Hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		token, refreshToken, err := tokens.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *&user.User_ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		user.Token = &token
		user.Refresh_Token = &refreshToken

		// Save to database
		result := database.DB.Create(&user)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return 
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Successfully signed up",
			"user_id": user.User_ID,
		})
	}
}

// func Login() gin.HandlerFunc {
	
// }

// func ProductViewerAdmin() gin.HandlerFunc {
	
// }

// func SearchProduct() gin.HandlerFunc {
	
// }

// func SearchProductByQuery() gin.HandlerFunc {
	
// }