package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adam-hanna/sessions"
	"github.com/adam-hanna/sessions/user"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

type Config struct {
	RedisAddress string
	AuthKey      string
	Https        bool
	Secure       bool
}

type Module struct {
	Service *sessions.Service
}

// 对一个request启动session
func (m *Module) StartSession(c *gin.Context, noCheckCsrf bool) {
	var (
		sessInfo = m.getUserSession(c)
	)

	c.Set("no_check_csrf", noCheckCsrf)

	if sessInfo == nil {
		// session未初始化
		mySession, err := m.NewSession()
		if err != nil {
			log.Printf("ERROR: NewSession: %v", err)
			return
		}
		JSONBytes, err := json.Marshal(mySession)
		if err != nil {
			log.Printf("ERROR: marhsalling json: %v", err)
			return
		}
		sessInfo, err = m.Service.IssueUserSession("", string(JSONBytes[:]), c.Writer)
		if err != nil {
			log.Printf("ERROR: issuing user session: %v", err)
			return
		}

		c.Request.Header.Set("X-CSRF-Token", mySession.CSRF)

		m.UpdateSessionInfo(c, sessInfo, &mySession)

	} else {
		m.Service.ExtendUserSession(sessInfo, c.Request, c.Writer)
	}

}

// 从session获取值
func (m *Module) GetValue(c *gin.Context, key string) (v interface{}, exist bool) {

	myJSON, sessErr := m.GetSessionJSON(c)
	if sessErr != nil {
		return
	}

	if myJSON == nil {
		return
	}

	v, exist = myJSON.Data[key]
	return
}

// 按字符串格式读取一个session值
func (m *Module) GetValueString(c *gin.Context, key string) string {

	v, exist := m.GetValue(c, key)
	if exist {
		return fmt.Sprint(v)
	}

	return ""
}

// 写入session
func (m *Module) SetValue(c *gin.Context, key string, v interface{}) {

	sessJson, sessErr := m.GetSessionJSON(c)
	if sessErr != nil || sessJson == nil {
		return
	}

	sessJson.Data[key] = v

	m.UpdateSessionJSON(c, sessJson)

}

// 移除值
func (m *Module) DelValue(c *gin.Context, key string) {

	sessJson, sessErr := m.GetSessionJSON(c)
	if sessErr != nil || sessJson == nil {
		return
	}

	delete(sessJson.Data, key)

	m.UpdateSessionJSON(c, sessJson)

}

// 清空session
func (m *Module) ClearSession(c *gin.Context) {
	var (
		sessInfo = m.getUserSession(c)
	)
	if sessInfo == nil {
		return
	}

	m.Service.ClearUserSession(sessInfo, c.Writer)

	// need to clear the csrf cookie, too
	aLongTimeAgo := time.Now().Add(-1000 * time.Hour)
	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    "",
		Expires:  aLongTimeAgo,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // note: can't use secure cookies in development
	}
	http.SetCookie(c.Writer, &csrfCookie)

}

func (m *Module) UpdateSessionJSON(c *gin.Context, sessJson *SessionJSON) {
	newJSON, jsonErr := json.Marshal(sessJson)
	if jsonErr != nil {
		return
	}

	// 更新
	sessInfo := m.getUserSession(c)
	if sessInfo == nil {
		return
	}

	sessInfo.JSON = string(newJSON)
	m.UpdateSessionInfo(c, sessInfo, sessJson)

}

func (m *Module) UpdateSessionInfo(c *gin.Context, sessInfo *user.Session, sessJson *SessionJSON) {

	c.Set("session", sessInfo)

	if err := m.Service.ExtendUserSession(sessInfo, c.Request, c.Writer); err != nil {
		log.Printf("ERROR: extending user session: %v", err)
		return
	}

	// need to extend the csrf cookie, too
	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    sessJson.CSRF,
		Expires:  sessInfo.ExpiresAt,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // note: can't use secure cookies in development
	}
	http.SetCookie(c.Writer, &csrfCookie)
}

func (m *Module) GetSessionJSON(c *gin.Context) (*SessionJSON, error) {
	var (
		sessInfo = m.getUserSession(c)
		myJSON   = SessionJSON{}
	)

	if err := json.Unmarshal([]byte(sessInfo.JSON), &myJSON); err != nil {
		log.Printf("ERROR: unmarshalling json: %v", err)
		return nil, err
	}

	if c.GetBool("no_check_csrf") {
		return &myJSON, nil
	}

	// 跨站请求伪造检查
	csrf := c.GetHeader("X-CSRF-Token")
	if csrf != myJSON.CSRF {
		csrfErr := errors.New("Unauthorized! CSRF token doesn't match user session")
		log.Printf("ERROR: %v", csrfErr)
		return nil, csrfErr
	}

	return &myJSON, nil
}

func (m *Module) getUserSession(c *gin.Context) *user.Session {

	v, exist := c.Get("session")

	if exist {
		return v.(*user.Session)
	}

	sessInfo, err := m.Service.GetUserSession(c.Request)
	if err != nil {
		log.Printf("ERROR: fetching user session: %v", err)
		return nil
	}

	return sessInfo
}

func (Module) NewSession() (SessionJSON, error) {

	var (
		s = SessionJSON{
			Data: map[string]interface{}{},
		}
		err error
	)

	s.CSRF, err = generateKey()

	return s, err
}

func generateKey() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
