package controllers

import (
	"net/http"

	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func ProductViewerAdmin(app *Application) gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		var product models.Product

		if err := ctx.BindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		err := app.ProductRepo.Create(&product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Product add successfully",
			"product_id": product.ID,
		})
	}
}

func SearchProduct(app *Application) gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		products, err := app.ProductRepo.FindAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", products)
	}
}

func SearchProductByQuery(app *Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("name")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return 
		}

		products, err := app.ProductRepo.FindByName(query)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to search product"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", products)
	}
}