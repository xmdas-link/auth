package wechat_qrcode

import "time"

type WechatUser struct {
	ID           uint32 `gorm:"primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	OpenId       string `gorm:"size:255;index:idx_open_id"`
	UnionId      string `gorm:"size:255;index:idx_union_id"`
	Name         string `gorm:"size:64;"`
	Token        string `gorm:"size:255"`
	TokenExpired *time.Time
	Role         string `gorm:"size:64"`
	Active       int32  `gorm:"default:1"`
}

func (u *WechatUser) IsActive() bool {
	return u.Active == 1
}
