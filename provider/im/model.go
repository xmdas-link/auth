package im

import "time"

type ImUser struct {
	ID           uint32 `gorm:"primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ImID         string `gorm:"size:255;index:idx_im_id"`
	Name         string `gorm:"size:64;"`
	LoginName    string `gorm:"size:64;"`
	Token        string `gorm:"size:255"`
	TokenExpired *time.Time
	Role         string `gorm:"size:64"`
	Active       int32  `gorm:"default:1"`
}

func (u *ImUser) IsActive() bool {
	return u.Active == 1
}
