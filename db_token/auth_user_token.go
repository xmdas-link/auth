package db_token

import "time"

type AuthUserToken struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Uid       string `gorm:"size:32;index:ind_uid_provider"`
	Provider  string `gorm:"size:32;index:ind_uid_provider"`
	Token     string `gorm:"size:32;index:idx_token"`
	TokenInfo string `gorm:"type:TEXT"`
	IP        string `gorm:"size:32"`
	ExpiredAt int64  `gorm:"index:idx_expired_at"`
	CreatedAt time.Time
}
