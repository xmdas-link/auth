package im

import (
	"github.com/xmdas-link/auth/user_store"
	"time"
)

type ImUser struct {
	ID           uint32 `gorm:"primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UID          string `gorm:"index:idx_uid"`
	ImID         string `gorm:"size:255;index:idx_im_id"`
	Name         string `gorm:"size:64;"`
	LoginName    string `gorm:"size:64;"`
	Token        string `gorm:"size:255"`
	TokenExpired *time.Time

	UserBase user_store.UserInterface `gorm:"-"`
}
