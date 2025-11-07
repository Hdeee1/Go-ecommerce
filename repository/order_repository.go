package repository

import (
	"github.com/Hdeee1/go-ecommerce/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindCartByUserID(userID uint) (models.Order, error) {
	var cartOrder models.Order
	err := r.db.Where("user_id = ? AND price = ?", userID).Preload("OrderItems").First(&cartOrder).Error

	return cartOrder, err
}

func (r *OrderRepository) CrateCart(cart *models.Order) error {
	return r.db.Create(cart).Error
}

func (r *OrderRepository) FindOrderItem(orderID uint, productID uint) (models.OrderItem, error) {
	var item models.OrderItem
	err := r.db.Where("order_id = ? AND product_id = ?", orderID, productID).First(&item).Error
	
	return item, err
}

func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *OrderRepository) UpdateOrderItem(item *models.OrderItem) error {
	return r.db.Save(item).Error
}

func (r *OrderRepository) DeleteOrderItem(orderID uint, productID uint) (int64, error) {
	result := r.db.Where("order_id = ? AND product_id = ?", orderID, productID).Delete(&models.OrderItem{})
	return result.RowsAffected, result.Error
}