package user_store

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint32 `gorm:"primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:64;"`
	Role      string `gorm:"size:64"`
	Active    int32  `gorm:"default:1"`
	HeaderUrl string `gorm:"size:255"`
}

func (u *User) IsActive() bool {
	return u.Active == 1
}

func (u *User) GetId() string {
	return fmt.Sprint(u.ID)
}

func (u *User) GetRole() string {
	return u.Role
}

func (u *User) GetMapData() map[string]string {
	return map[string]string{
		"id":         fmt.Sprint(u.ID),
		"created_at": fmt.Sprint(u.CreatedAt.Unix()),
		"updated_at": fmt.Sprint(u.UpdatedAt.Unix()),
		"name":       u.Name,
		"role":       u.Role,
		"active":     fmt.Sprint(u.Active),
		"header_url": u.HeaderUrl,
	}
}
