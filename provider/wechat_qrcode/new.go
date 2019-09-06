package wechat_qrcode

func New(cfg *OAuthConfig) *Provider {

	if cfg.AuthBaseUrl == "" {
		cfg.AuthBaseUrl = "https://open.weixin.qq.com"
	}

	var (
		p = &Provider{
			Name: Name,
			Client: &Client{
				cfg,
			},
		}
	)

	return p
}
