package repositories

import (
	"strings"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func (u *UserStore) Create(user *models.User) error {
	return u.db.Create(&user).Error
}

func (u *UserStore) CreateToken(userToken *models.UserToken) error {
	return u.db.Create(&userToken).Error
}

func (u *UserStore) GetAll() ([]*models.User, error) {
	var users []*models.User
	err := u.db.Find(&users).Error
	return users, err
}

func (u *UserStore) GetOne(id uint64) (*models.User, error) {
	var user *models.User
	err := u.db.First(&user, "id = ?", id).Error
	return user, err
}

func (u *UserStore) GetByEmail(email string) (*models.User, error) {
	var user *models.User
	result := u.db.Find(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	return user, nil
}

func (u *UserStore) GetByUsername(username string) (*models.User, error) {
	var user *models.User
	result := u.db.Find(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	return user, nil
}

func (u *UserStore) GetByGoogleID(googleID string) (*models.User, error) {
	var user *models.User
	result := u.db.Find(&user, "google_id = ?", googleID)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	return user, nil
}

func (u *UserStore) GetByRefreshToken(userId uint64, token string) (*models.User, error) {
	var user *models.User
	err := u.db.
		Joins("JOIN user_tokens ut ON ut.user_id = users.id").
		Where("ut.user_id = ? AND ut.refresh_token = ?", userId, token).
		Take(&user).Error
	return user, err
}

func (u *UserStore) GetByEmailOrUsername(credential string) (*models.User, error) {
	if strings.Contains(credential, "@") {
		return u.GetByEmail(credential)
	}
	return u.GetByUsername(credential)
}

func (u *UserStore) Update(user *models.User) (*models.User, error) {
	err := u.db.Model(&user).Updates(user).Error
	return user, err
}

func (u *UserStore) UpdateToken(userToken *models.UserToken) error {
	return u.db.Save(userToken).Error
}

func (u *UserStore) DeleteRefreshToken(userId uint64, token string) error {
	return u.db.
		Model(&models.UserToken{}).
		Where("user_id = ? AND refresh_token = ?", userId, token).
		Update("refresh_token", "").
		Error
}

func (u *UserStore) Delete(id uint64) error {
	return nil
}
