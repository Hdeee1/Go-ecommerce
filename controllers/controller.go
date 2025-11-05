package controllers

import (
	"net/http"

	"github.com/Hdeee1/go-ecommerce/database"
	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/tokens"
	"github.com/Hdeee1/go-ecommerce/utils"
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

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	if err != nil {
		return  false, "Login or Password is incorrect"
	}

	return true, ""
}

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

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.DB.Where("email = ?", user.Email).First(&foundUser).Error
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Login or password is incorrect"})
			return 
		}

		// verify password
		isValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return 
		}

		// generate token 
		token, refreshToken, err := tokens.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *&foundUser.User_ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		database.DB.Model(&foundUser).Updates(models.User{
			Token: &token,
			Refresh_Token: &refreshToken,
		})

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message":		 	"Successfully logged in",
			"token":			token,
			"refresh_token": 	refreshToken,
			"user_id": 			foundUser.User_ID,
		})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		var product models.Product

		if err := ctx.BindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		result := database.DB.Create(&product)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Product add successfully",
			"product_id": product.ID,
		})
	}
}

func SearchProduct() gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		var products []models.Product

		result := database.DB.Find(&products)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", products)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var products []models.Product

		query := ctx.Query("name")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return 
		}

		result := database.DB.Where("product_name LIKE ?", "%"+query+"%").Find(&products)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to search product"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", products)
	}
}