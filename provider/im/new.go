package im

import "golang.org/x/oauth2"

func New(cfg *OAuthConfig) *Provider {
	var (
		p = &Provider{
			Name:       Name,
			BaseUrl:    cfg.MattermostUrl,
			UrlVersion: "/api/v4",
		}
	)

	p.oauthCfg = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.Secret,
		RedirectURL:  cfg.CallbackUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.MattermostUrl + "/oauth/authorize",
			TokenURL: cfg.MattermostUrl + "/oauth/access_token",
		},
		Scopes: []string{},
	}
	return p
}
