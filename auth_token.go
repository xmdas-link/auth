package auth

type AuthToken interface {

	// 新Token
	NewToken(user map[string]string) (token string, expiredAt int64, err error)

	// 清除Token
	ClearToken(token string) error

	// 清除用户的Token
	ClearTokenOfUser(uid string, provider string) error

	// 查找token
	FindToken(token string) (user map[string]string)
}
