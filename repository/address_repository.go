package repository

import (
	"github.com/Hdeee1/go-ecommerce/models"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(address *models.Address) error {
	return r.db.Create(address).Error
}

func (r *AddressRepository) Update(address *models.Address) error {
	return r.db.Save(address).Error
}

func (r *AddressRepository) FindByUserAndType(userID uint, addrType string) (*models.Address, error) {
	var address models.Address
	err := r.db.Where("user_id = ? AND type = ?", userID, addrType).First(&address).Error

	return &address, err
}

func (r *AddressRepository) FindByID(id uint) (*models.Address, error) {
	var address models.Address
	err := r.db.Where("user_id = ? AND address_id", id, address.ID).First(&address, id).Error

	return &address, err
}

func (r *AddressRepository) Delete(id uint) error {
	err := r.db.Where("address_id = ?", id).Delete(&models.Address{})

	return err.Error
}