package repositories

import (
	"gorm.io/gorm"
)

type Repositories struct {
	Users *UserStore
}

func NewRepository(db *gorm.DB) Repositories {
	return Repositories{
		Users: &UserStore{db},
	}
}
