package im

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	URLVersion string
	Token      *oauth2.Token
	Im         *http.Client
}

func (c *Client) GetMe() (user ImUserData, err error) {

	rsp, httpErr := c.Im.Get(c.BaseURL + c.URLVersion + "/users/me")
	if httpErr != nil {
		err = fmt.Errorf("请求IM API发生错误：%v", httpErr)
		return
	}

	body, errRead := ioutil.ReadAll(rsp.Body)
	if errRead != nil {
		err = fmt.Errorf("读取IM API返回结果失败：%v", errRead)
		return
	}
	defer rsp.Body.Close()

	log.Print(string(body))

	if rsp.StatusCode != 200 {
		errorData := ErrorData{}
		if jsonErr := json.Unmarshal(body, &errorData); jsonErr != nil {
			err = fmt.Errorf("转换IM API返回的错误信息失败：%v", jsonErr)
		} else {
			err = errors.New(errorData.Message)
		}
	} else {
		if jsonErr := json.Unmarshal(body, &user); jsonErr != nil {
			err = fmt.Errorf("转换IM API返回的用户信息失败：%v", jsonErr)
		}
	}

	return
}

func (c *Client) GetAccessToken() (token string, expired *time.Time) {
	if c.Token != nil {
		token = c.Token.AccessToken
		expired = &c.Token.Expiry
	}
	return
}

type ErrorData struct {
	StatusCode int    `json:"status_code"`
	ID         string `json:"id"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
}

type ImUserData struct {
	ID            string `json:"id"`
	CreateAt      int    `json:"create_at"`
	UpdateAt      int    `json:"update_at"`
	DeleteAt      int    `json:"delete_at"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AuthService   string `json:"auth_service"`
	Roles         string `json:"roles"`
	Locale        string `json:"locale"`
	NotifyProps   struct {
		Email        string `json:"email"`
		Push         string `json:"push"`
		Desktop      string `json:"desktop"`
		DesktopSound string `json:"desktop_sound"`
		MentionKeys  string `json:"mention_keys"`
		Channel      string `json:"channel"`
		FirstName    string `json:"first_name"`
	} `json:"notify_props"`
	Props struct {
	} `json:"props"`
	LastPasswordUpdate int  `json:"last_password_update"`
	LastPictureUpdate  int  `json:"last_picture_update"`
	FailedAttempts     int  `json:"failed_attempts"`
	MfaActive          bool `json:"mfa_active"`
	Timezone           struct {
		UseAutomaticTimezone string `json:"useAutomaticTimezone"`
		ManualTimezone       string `json:"manualTimezone"`
		AutomaticTimezone    string `json:"automaticTimezone"`
	} `json:"timezone"`
	TermsOfServiceID       string `json:"terms_of_service_id"`
	TermsOfServiceCreateAt int    `json:"terms_of_service_create_at"`
}
