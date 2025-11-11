package controllers

import (
	"net/http"
	"strconv"

	"github.com/Hdeee1/go-ecommerce/models"
	"github.com/Hdeee1/go-ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func (app *Application) AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		var address models.Address

		if err := ctx.BindJSON(&address); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		address.UserID = user.ID

		if err := app.AddressRepo.Create(&address); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusCreated, "Address added successfully", gin.H{
			"address_id": address.ID,
		})
	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := ctx.Get("user_id")
		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		
		var newAddressData models.Address
		if err := ctx.BindJSON(&newAddressData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := Validate.Struct(newAddressData); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		addressType := "Home"
		existingAddress, err := app.AddressRepo.FindByUserAndType(user.ID, addressType)

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				newAddressData.UserID = user.ID
				newAddressData.Type = &addressType

				if createErr := app.AddressRepo.Create(&newAddressData); createErr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create home address"})
					return
				}
				utils.SuccessResponse(ctx, http.StatusCreated, "Home address created", nil)
				return

			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find address"})
				return
			}
		}

		existingAddress.House = newAddressData.House
		existingAddress.Street = newAddressData.Street
		existingAddress.City = newAddressData.City
		existingAddress.Pincode = newAddressData.Pincode

		if updateErr := app.AddressRepo.Update(existingAddress); updateErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update home address"})
			return
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Home address updated", nil)
	}
}

func (app *Application) EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := ctx.Get("user_id")
		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		
		var newAddressData models.Address
		if err := ctx.BindJSON(&newAddressData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := Validate.Struct(newAddressData); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		addressType := "Work"
		existingAddress, err := app.AddressRepo.FindByUserAndType(user.ID, addressType)

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				newAddressData.UserID = user.ID
				newAddressData.Type = &addressType

				if createErr := app.AddressRepo.Create(&newAddressData); createErr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Work address"})
					return
				}
				utils.SuccessResponse(ctx, http.StatusCreated, "Work address created", nil)
				return

			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find address"})
				return
			}
		}

		existingAddress.House = newAddressData.House
		existingAddress.Street = newAddressData.Street
		existingAddress.City = newAddressData.City
		existingAddress.Pincode = newAddressData.Pincode

		if updateErr := app.AddressRepo.Update(existingAddress); updateErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Work address"})
			return
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Work address updated", nil)
	}
}

func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := ctx.Get("user_id")
		user, err := app.UserRepo.FindByID(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return 
		}

		addrIDStr := ctx.Param("address_id")

		addressID, err := strconv.ParseUint(addrIDStr, 10, 23)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address id"})
			return 
		}

		addressToDelete, err := app.AddressRepo.FindByID(uint(addressID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Address id not found"})
				return 
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find address"})
			return 
		}

		if addressToDelete.UserID != user.ID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorize to delete this address"})
			return 
		}

		if err := app.AddressRepo.Delete(uint(addressID)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
			return 
		}

		utils.SuccessResponse(ctx, http.StatusOK, "Address deleted successfully", nil)
	}
}