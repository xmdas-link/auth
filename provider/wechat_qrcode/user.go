package wechat_qrcode

import "fmt"

type UserData struct {
	Provider  string
	Token     string
	ExpiredAt int64
	*WechatUser
}

func (u *UserData) GetID() string {
	return fmt.Sprint(u.ID)
}

func (u *UserData) GetProvider() string {
	return u.Provider
}
func (u *UserData) GetRole() string {
	return u.Role
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
		"role":     u.Role,
		"user":     u.OpenId,
		"union_id": u.UnionId,
		"active":   fmt.Sprint(u.Active),
	}
}
