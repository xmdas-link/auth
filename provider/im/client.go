package im

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*OAuthConfig
	/*BaseURL    string
	URLVersion string
	Token      *oauth2.Token
	Im         *http.Client*/
}

func (c *Client) AuthCodeURL(state string) string {

	// 拼接登录地址
	var (
		buf bytes.Buffer
		v   = url.Values{}
	)

	v.Set("response_type", "code")
	v.Set("client_id", c.ClientID)
	v.Set("redirect_uri", c.CallbackUrl)
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	if state != "" {
		v.Set("state", state)
	}

	buf.WriteString(c.MattermostUrl + "/oauth/authorize?")
	buf.WriteString(v.Encode())
	return buf.String()
}

func (c *Client) GetMe(token *AccessToken) (user ImUserData, err error) {

	var (
		client = &http.Client{}
	)
	req, reqErr := http.NewRequest("GET", c.MattermostUrl+c.ApiVersion+"/users/me", strings.NewReader(""))
	if reqErr != nil {
		err = fmt.Errorf("请求IM API发生错误：%v", reqErr)
		return
	}

	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)

	rsp, httpErr := client.Do(req)
	if httpErr != nil {
		err = fmt.Errorf("请求IM API发生错误：%v", httpErr)
		return
	}
	defer rsp.Body.Close()

	body, errRead := ioutil.ReadAll(rsp.Body)
	if errRead != nil {
		err = fmt.Errorf("读取IM API返回结果失败：%v", errRead)
		return
	}

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

func (c *Client) GetAccessToken(code string) (*AccessToken, error) {
	var (
		v        = url.Values{}
		tokenURL = c.MattermostUrl + "/oauth/access_token"
		token    = AccessToken{}
	)
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("redirect_uri", c.CallbackUrl)
	v.Set("client_id", c.ClientID)
	v.Set("client_secret", c.Secret)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Print(string(body))
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // at least PayPal returns string, while most return number
	Expires      int64  `json:"expires"`    // broken Facebook spelling of expires_in
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
