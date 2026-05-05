package order

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db}
}

func (r *Repository) CreateOrder(order *Order) error {
    return r.db.Create(order).Error
}