package im

import "fmt"

type UserData struct {
	Provider  string
	Token     string
	ExpiredAt int64
	*ImUser
}

func (u *UserData) GetID() string {
	if u.ImUser.UserBase != nil {
		return u.ImUser.UserBase.GetId()
	}
	return ""
}

func (u *UserData) GetProvider() string {
	return u.Provider
}
func (u *UserData) GetRole() string {
	if u.ImUser.UserBase != nil {
		return u.UserBase.GetRole()
	}
	return ""
}

func (u *UserData) GetToken() string {
	return u.Token
}

func (u *UserData) GetExpired() int64 {
	return u.ExpiredAt
}

func (u *UserData) SetToken(v string, expire int64) {
	u.Token = v
	u.ExpiredAt = expire
}

func (u *UserData) GetMapData() map[string]string {
	return map[string]string{
		"id":       fmt.Sprint(u.ID),
		"provider": u.Provider,
		"name":     u.Name,
		"role":     u.GetRole(),
		"user":     u.LoginName,
	}
}
