package order

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateOrder(order *Order) error {
	return r.DB.Create(order).Error
}

func (r *Repository) FindOrdersByUser(userID uint) ([]Order, error) {
	var orders []Order

	err := r.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error

	return orders, err
}

func (r *Repository) FindOrderByID(id uint) (*Order, error) {
	var order Order

	err := r.DB.
		Preload("Items").
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) SaveOrder(order *Order) error {
	return r.DB.Save(order).Error
}
func (r *Repository) WithTx(tx *gorm.DB) *Repository {
	return &Repository{
		DB: tx,
	}
}