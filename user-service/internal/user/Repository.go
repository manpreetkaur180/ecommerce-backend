package user

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}
func (r *Repository) CreateAddress(address *Address) error {
	return r.DB.Create(address).Error
}

func (r *Repository) RemoveDefaultAddresses(userID uint) error {
	return r.DB.Model(&Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

func (r *Repository) FindAddressesByUserID(userID uint) ([]Address, error) {
	var addresses []Address

	err := r.DB.
		Where("user_id = ?", userID).
		Order("is_default desc, created_at desc").
		Find(&addresses).Error

	return addresses, err
}

func (r *Repository) CountUserAddresses(userID uint) (int64, error) {
	var count int64

	err := r.DB.Model(&Address{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count, err
}

func (r *Repository) UnsetDefaultAddresses(userID uint) error {
	return r.DB.Model(&Address{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}


func (r *Repository) CreateUser(u *User) error {
	return r.DB.Create(u).Error
}

func (r *Repository) FirstUserByEmail(email string) (*User, error) {
	var u User
	err := r.DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FirstUserByPhone(phone string) (*User, error) {
	var u User
	err := r.DB.Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FirstUserByID(id uint) (*User, error) {
	var u User
	err := r.DB.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FirstUserByEmailOrPhone(email, phone string) (*User, error) {
	var u User
	var err error
	if email != "" {
		err = r.DB.Where("email = ?", email).First(&u).Error
	} else {
		err = r.DB.Where("phone = ?", phone).First(&u).Error
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FirstUserForLogin(email, phone string) (*User, error) {
	var u User
	query := r.DB
	if email != "" {
		query = query.Where("email = ?", email)
	} else {
		query = query.Where("phone = ?", phone)
	}
	if err := query.First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateEmailVerification(v *EmailVerification) error {
	return r.DB.Create(v).Error
}

func (r *Repository) FirstEmailVerificationByTokenHash(tokenHash string) (*EmailVerification, error) {
	var v EmailVerification
	err := r.DB.Where("token_hash = ?", tokenHash).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) SaveUser(u *User) error {
	return r.DB.Save(u).Error
}

func (r *Repository) SaveEmailVerification(v *EmailVerification) error {
	return r.DB.Save(v).Error
}

func (r *Repository) CreatePasswordUpdate(p *PasswordUpdate) error {
	return r.DB.Create(p).Error
}

func (r *Repository) FirstPasswordUpdateByTokenHash(tokenHash string) (*PasswordUpdate, error) {
	var p PasswordUpdate
	err := r.DB.Where("token_hash = ?", tokenHash).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) DeletePasswordUpdate(p *PasswordUpdate) error {
	return r.DB.Delete(p).Error
}
