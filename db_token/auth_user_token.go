package db_token

import "time"

var (
	TokenTableName = "op_auth_user_token"
)

type AuthUserToken struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Uid       string `gorm:"size:32;index:ind_uid_provider"`
	Provider  string `gorm:"size:32;index:ind_uid_provider"`
	Token     string `gorm:"size:32;index:idx_token"`
	TokenInfo string `gorm:"type:TEXT"`
	IP        string `gorm:"size:32"`
	ExpiredAt int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AuthUserToken) TableName() string {
	return TokenTableName
}
