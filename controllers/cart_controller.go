package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productIDStr := ctx.Query("product_id")
		if productIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
			return
		}

		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return 
		}

		userID, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return 
		}

		product, err := app.ProductRepo.FindByID(uint(productID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return 
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return 
		}

		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return 
		}

		cartOrder, err := app.OrderRepo.FindCartByUserID(user.ID)
		if err == gorm.ErrRecordNotFound {
			cartOrder = models.Order{
				UserID: user.ID,
				Price: 0,
			}
			if err := app.OrderRepo.CreateCart(&cartOrder); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return 
			}
		}

		existingItem, err := app.OrderRepo.FindOrderItem(cartOrder.ID, uint(productID))
		if err == gorm.ErrRecordNotFound {
			quantity := uint(1)
			if qtyStr := ctx.Query("quantity"); qtyStr != "" {
				if qty, err := strconv.ParseUint(qtyStr, 10, 32); err == nil {
					quantity = uint(qty)
				}
			}

			orderItem := models.OrderItem{
				OrderID:      cartOrder.ID,
				ProductID:    uint(productID),
				Product_Name: product.Product_Name,
				Price:        *product.Price,
				Rating:       product.Rating,
				Image:        product.Image,
				Quantity:     quantity,
			}

			if err := app.OrderRepo.CreateOrderItem(&orderItem); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
				return 
			}
		} else {
			existingItem.Quantity += 1
			if err := app.OrderRepo.UpdateOrderItem(&existingItem); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the cart item"})
				return
			}
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message": "Product added to cart successfully",
			"cart_id": cartOrder.ID,
		})
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productIDStr := ctx.Query("product_id")
		if productIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
			return 
		}

		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
			return 
		}

		userID, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return 
		}
		
		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		cartOrder, err := app.OrderRepo.FindCartByUserID(user.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Cart is empty"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return 
		}

		rowsAffected, err := app.OrderRepo.DeleteOrderItem(cartOrder.ID, uint(productID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
			return 
		}

		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message": "Item removed from cart successfully",
		})
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
			return 
		}

		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user_id not found"})
			return 
		}

		cartOrder, err := app.OrderRepo.FindCartByUserID(user.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "cart is empty"})
				return 
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return 
		}
		
		if len(cartOrder.OrderItems) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
			return 
		}

		var totalPrice uint64 = 0
		for _, item := range cartOrder.OrderItems {
			totalPrice += item.Price * uint64(item.Quantity)
		}

		paymentMethod := ctx.Query("payment_method")
		var payment models.Payment
		if paymentMethod == "digital" {
			payment = models.Payment{
				Digital: true,
				COD: false,
			}
		} else {
			payment = models.Payment{
				Digital: false,
				COD: true,
			}
		}

		cartOrder.Price = int(totalPrice)
		cartOrder.Order_At = time.Now()
		cartOrder.Payment_Method = payment

		if err := app.OrderRepo.UpdateOrder(&cartOrder); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to place order"})
			return
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message": "Order placed successfully",
			"order_id": cartOrder.ID,
			"total_price": totalPrice,
			"items_count": len(cartOrder.OrderItems),
			"payment": paymentMethod,
		}) 
	}	
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type InstantBuyRequest struct {
			ProductID		uint	`json:"product_id" binding:"required"`
			Quantity		uint	`json:"quantity" binding:"required,min=1"`
			PaymentMethod	string	`json:"payment_method" binding:"required"`
		}

		var req InstantBuyRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		userID, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authenticated"})
			return 
		}

		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return 
		}

		product, err := app.ProductRepo.FindByID(req.ProductID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return 
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return 
		}

		totalPrice := *product.Price * uint64(req.Quantity)

		var payment models.Payment
		if req.PaymentMethod == "digital" {
			payment = models.Payment{Digital: true, COD: false}
		} else {
			payment = models.Payment{Digital: false, COD: true}
		}

		order := models.Order{
			UserID: user.ID,
			Price: int(totalPrice),
			Order_At: time.Now(),
			Payment_Method: payment,
		}

		if err := app.OrderRepo.CreateCart(&order); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return 
		}

		orderItem := models.OrderItem{
			OrderID:      order.ID,
			ProductID:    req.ProductID,
			Product_Name: product.Product_Name,
			Price:        *product.Price,
			Rating:       product.Rating,
			Image:        product.Image,
			Quantity:     req.Quantity,
		}

		if err := app.OrderRepo.CreateOrderItem(&orderItem); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Success", gin.H{
			"message": "Order placed successfully",
			"order_id": order.ID,
			"total_price": totalPrice,
			"quantity": req.Quantity,
			"payment": req.PaymentMethod,
		})
	}
}