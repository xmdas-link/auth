package wechat_qrcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*OAuthConfig
}

func (c *Client) AuthCodeData(state string) map[string]string {
	data := map[string]string{
		"appid":         c.ClientID,
		"response_type": "code",
		"redirect_uri":  c.CallbackUrl,
		"state":         state,
		"scope":         strings.Join(c.Scopes, ","),
	}

	return data
}

func (c *Client) AuthCodeURL(state string) string {
	// 拼接微信登录地址
	var (
		buf bytes.Buffer
		v   = url.Values{}
	)
	buf.WriteString(c.AuthBaseUrl + "/connect/qrconnect?")

	for key, value := range c.AuthCodeData(state) {
		v.Set(key, value)
	}

	buf.WriteString(v.Encode())
	buf.WriteString("#wechat_redirect")

	return buf.String()
}

func (c *Client) GetAccessToken(code string) (*AccessToken, error) {
	var (
		path  = "https://api.weixin.qq.com/sns/oauth2/access_token"
		v     = url.Values{}
		token = AccessToken{}
	)

	v.Set("appid", c.ClientID)
	v.Set("secret", c.Secret)
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")

	if err := c.DoGet(path, v, &token); err != nil {
		return nil, fmt.Errorf("请求微信API发生错误：%v", err)
	}
	return &token, nil
}

func (c *Client) GetWechatUser(token *AccessToken) (*WechatUserInfo, error) {
	var (
		path = "https://api.weixin.qq.com/sns/userinfo"
		v    = url.Values{}
		user = WechatUserInfo{}
	)
	v.Set("access_token", token.AccessToken)
	v.Set("openid", token.Openid)

	if err := c.DoGet(path, v, &user); err != nil {
		return nil, fmt.Errorf("请求微信API发生错误：%v", err)
	}
	return &user, nil
}

func (c *Client) DoGet(path string, v url.Values, data interface{}) error {

	resp, err := http.Get(path + "?" + v.Encode())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Print(string(body))

	result := make(map[string]interface{})
	if jsonErr := json.Unmarshal(body, &result); jsonErr != nil {
		return jsonErr
	}

	if errMsg, ok := result["errmsg"]; ok {
		return fmt.Errorf("%v，%v", result["errcode"], errMsg)
	}

	jsonErr := json.Unmarshal(body, data)
	return jsonErr

}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

type WechatUserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}
