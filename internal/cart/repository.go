package cart

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db}
}

// get cart by user
func (r *Repository) GetCartByUser(userID uint) (*Cart, error) {
    var cart Cart
    err := r.db.Preload("Items").
        Where("user_id = ?", userID).
        First(&cart).Error

    return &cart, err
}

// create cart
func (r *Repository) CreateCart(cart *Cart) error {
    return r.db.Create(cart).Error
}

// add item
func (r *Repository) AddItem(item *CartItem) error {
    return r.db.Create(item).Error
}

// update item
func (r *Repository) UpdateItem(item *CartItem) error {
    return r.db.Save(item).Error
}

// delete item
func (r *Repository) RemoveItem(cartID, productID uint) error {
    return r.db.
        Where("cart_id = ? AND product_id = ?", cartID, productID).
        Delete(&CartItem{}).Error
}