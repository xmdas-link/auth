package user_store

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

type UserStoreInterface interface {
	Get(id string, c *gin.Context) (UserInterface, error)
	New(data map[string]string, c *gin.Context) (user UserInterface, id string, err error)
}

type UserInterface interface {
	GetId() string
	GetRole() string
	IsActive() bool
	GetMapData() map[string]string
}

func New(db *gorm.DB) *UserStore {
	if db == nil {
		return nil
	}

	db.AutoMigrate(&User{})

	return &UserStore{
		DB: db,
	}
}

type UserStore struct {
	DB *gorm.DB
}

func (u *UserStore) Get(id string, c *gin.Context) (UserInterface, error) {

	var (
		tx   = u.DB
		data = User{}
	)

	err := tx.First(&data, "id = ?", id).Error
	return &data, err
}

func (u *UserStore) New(info map[string]string, c *gin.Context) (user UserInterface, id string, err error) {

	var (
		tx   = u.DB
		data = User{}
	)

	// 从info中读取需要的字段
	if name, exist := info["name"]; exist {
		data.Name = name
	}

	if role, exist := info["role"]; exist {
		data.Role = role
	}

	if headerUrl, exist := info["header_url"]; exist {
		data.HeaderUrl = headerUrl
	}

	if active, exist := info["active"]; exist {
		a, _ := strconv.Atoi(active)
		data.Active = int32(a)
	}

	err = tx.Create(&data).Error
	if err == nil {
		user = &data
		id = fmt.Sprint(data.ID)
	}

	return

}
