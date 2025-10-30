package controllers

import (
	"net/http"
	"strconv"

	"github.com/Hdeee1/go-ecommerce/database"
	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	ProductDB 	*gorm.DB
	UserDB		*gorm.DB
}

func NewApplication(productDB, userDB *gorm.DB) *Application {
	return &Application{
		ProductDB: productDB,
		UserDB: userDB,
	}
}

func (app *Application) AddToCart() func() {
	return func(ctx *gin.Context) {
		productIDStr := ctx.Query("product_id")
		if productIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
			return
		}

		productID, err := strconv.ParseInt(productIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return 
		}

		userID, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return 
		}

		var product models.Product
		if err := database.DB.First(&product, productID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return 
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return 
		}

		var user models.User
		if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return 
		}

		var cartOrder models.Order
		err = database.DB.Where("user_id = ? AND price = ?", userID, 0).
			Preload("OrderItems").
			First(&cartOrder).Error

		if err == gorm.ErrRecordNotFound {
			cartOrder = models.Order{
				UserID: user.ID,
				Price: 0,
			}
			if err := database.DB.Create(&cartOrder).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return 
			}
		}

		var existingItem models.OrderItem
		err = database.DB.Where("order_id = ?  AND product_id = ?", cartOrder.ID, productID).First(&existingItem).Error
		if err == gorm.ErrRecordNotFound {
			quantity := uint(1)
			if qtyStr := ctx.Query("quantity"); qtyStr != "" {
				if qty, err := strconv.ParseUint(qtyStr, 10, 32); err == nil {
					quantity = uint(qty)
				}
			}

			orderItem := models.OrderItem{
					OrderID: 		cartOrder.ID,
					ProductID: 		uint(productID),
					Product_Name: 	product.Product_Name,
					Price: 			*product.Price,
					Rating: 		product.Rating,
					Image: 			product.Image,
					Quantity: 		quantity,
			}

			if err := database.DB.Create(&orderItem).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
				return 
			}
		} else {
			existingItem.Quantity += 1
			if err := database.DB.Save(&existingItem).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the cart item"})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Product added to cart successfully",
			"cart_id": cartOrder.ID,
		})
	}
}

func (app *Application) RemoveItem() func() {
	return func() {
		
	}

	
}

func (app *Application) BuyFromCart() func() {
	return func() {
		
	}

	
}

func (app *Application) InstantBuy() func() {
	return func() {
		
	}

	
}