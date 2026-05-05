package product

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db}
}

func (r *Repository) Create(product *Product) error {
    return r.db.Create(product).Error
}

func (r *Repository) FindAll() ([]Product, error) {
    var products []Product
    err := r.db.Find(&products).Error
    return products, err
}