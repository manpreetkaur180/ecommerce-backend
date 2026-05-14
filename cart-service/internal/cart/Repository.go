package cart

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetCartByUserID(userID uint) (*Cart, error) {

	var cart Cart

	err := r.DB.
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *Repository) CreateCart(userID uint) (*Cart, error) {

	cart := Cart{
		UserID: userID,
	}

	if err := r.DB.Create(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *Repository) GetCartItem(
	cartID uint,
	productID uint,
) (*CartItem, error) {

	var item CartItem

	err := r.DB.
		Where(
			"cart_id = ? AND product_id = ?",
			cartID,
			productID,
		).
		First(&item).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) CreateCartItem(item *CartItem) error {
	return r.DB.Create(item).Error
}

func (r *Repository) UpdateCartItem(item *CartItem) error {
	return r.DB.Save(item).Error
}

func (r *Repository) DeleteCartItem(item *CartItem) error {
	return r.DB.Delete(item).Error
}

func (r *Repository) DeleteCartItemByProductID(
	cartID uint,
	productID uint,
) error {

	result := r.DB.
		Where(
			"cart_id = ? AND product_id = ?",
			cartID,
			productID,
		).
		Delete(&CartItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (r *Repository) GetCartItems(cartID uint) ([]CartItem, error) {

	var items []CartItem

	err := r.DB.
		Where("cart_id = ?", cartID).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) ClearCart(cartID uint) error {

	return r.DB.
		Where("cart_id = ?", cartID).
		Delete(&CartItem{}).Error
}