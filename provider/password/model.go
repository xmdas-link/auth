package password

import "time"

var (
	UserTableName = "op_user"
)

type User struct {
	ID        uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:64;"`
	LoginName string `gorm:"size:64;unique_index;not null"`
	Password  string `gorm:"size:255"`
	Role      string `gorm:"size:64"`
	Active    int32  `gorm:"default:1"`
}

func (u *User) IsActive() bool {
	return u.Active == 1
}

func (User) TableName() string {
	return UserTableName
}
