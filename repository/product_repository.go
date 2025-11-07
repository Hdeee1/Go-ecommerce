package repository

import (
	"github.com/Hdeee1/go-ecommerce/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error

	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error

	return &product, err
}

func (r *ProductRepository) FindByName(query string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("product_name LIKE ?", "%"+query+"%").Find(&products).Error

	return products, err
}