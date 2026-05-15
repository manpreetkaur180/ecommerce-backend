package product

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(p *Product) error {
	return r.DB.Create(p).Error
}

func (r *Repository) CountActive() (int64, error) {
	var total int64
	err := r.DB.Model(&Product{}).Where("is_active = ?", true).Count(&total).Error
	return total, err
}

func (r *Repository) FindActiveOrderedPaginated(limit, offset int) ([]Product, error) {
	var out []Product
	err := r.DB.Where("is_active = ?", true).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&out).Error
	return out, err
}

func (r *Repository) FirstActiveByID(id uint) (*Product, error) {
	var p Product
	err := r.DB.Where("id = ? AND is_active = ?", id, true).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) FirstBySellerAndProductID(sellerID, productID uint) (*Product, error) {
	var p Product
	err := r.DB.Where("id = ? AND seller_id = ?", productID, sellerID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) CountBySeller(sellerID uint) (int64, error) {
	var total int64
	err := r.DB.Model(&Product{}).Where("seller_id = ?", sellerID).Count(&total).Error
	return total, err
}

func (r *Repository) FindBySellerPaginated(sellerID uint, limit, offset int) ([]Product, error) {
	var out []Product
	err := r.DB.Where("seller_id = ?", sellerID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&out).Error
	return out, err
}

func (r *Repository) FindActiveByIDs(ids []uint) ([]Product, error) {
	var out []Product
	if len(ids) == 0 {
		return out, nil
	}
	err := r.DB.Where("id IN ? AND is_active = ?", ids, true).Find(&out).Error
	return out, err
}

func (r *Repository) Save(p *Product) error {
	return r.DB.Save(p).Error
}
func (r *Repository) FindProductByID(
	id uint,
) (*Product, error) {

	var product Product

	err := r.DB.First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *Repository) ReduceStock(
	productID uint,
	quantity int,
) error {

	result := r.DB.Exec(`
		UPDATE products
		SET stock = stock - ?
		WHERE id = ?
		AND stock >= ?
	`,
		quantity,
		productID,
		quantity,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrInvalidData
	}

	return nil
}