package controllers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/Hdeee1/go-ecommerce/tokens"
	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/utils"
	"github.com/gin-gonic/gin"
)

var Validate = validator.New()

func SignUp(app *Application) gin.HandlerFunc {
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

		_, err := app.UserRepo.FindByEmail(*user.Email)
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist"})
			return 
		}

		_, err = app.UserRepo.FindByPhone(*user.Phone)
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone already exist"})
			return
		}

		password := utils.HashPassword(*user.Password)
		user.Password = &password

		token, refreshToken, err := tokens.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		user.Token = &token
		user.Refresh_Token = &refreshToken

		err = app.UserRepo.Create(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return 
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Successfully signed up",
			"user_id": user.User_ID,
		})
	}
}

func Login(app *Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		foundUser, err := app.UserRepo.FindByEmail(*user.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Login or password is incorrect"})
			return 
		}

		isValid, msg := utils.VerifyPassword(*user.Password, *foundUser.Password)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return 
		}

		token, refreshToken, err := tokens.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		foundUser.Token = &token
		foundUser.Refresh_Token = &refreshToken
		app.UserRepo.Update(foundUser)

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message":		 	"Successfully logged in",
			"token":			token,
			"refresh_token": 	refreshToken,
			"user_id": 			foundUser.User_ID,
		})
	}
}