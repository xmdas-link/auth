package session

// SessionJSON is used for marshalling and unmarshalling custom session json information.
// We're using it as an opportunity to tie csrf strings to sessions to prevent csrf attacks
type SessionJSON struct {
	CSRF string                 `json:"csrf"` // 用来防止Cross-site request forgery攻击（多窗口浏览器的cookie安全问题）
	Data map[string]interface{} `json:"data"`
}
