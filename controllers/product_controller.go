package controllers

import (
	"net/http"

	"github.com/Hdeee1/go-ecommerce/database"
	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/utils"
	"github.com/gin-gonic/gin"
)



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