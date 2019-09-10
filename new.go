package auth

import (
	"errors"
)

func New(cfg Config) (*GinAuth, error) {

	var (
		ga = &GinAuth{}
	)

	// 默认配置
	if cfg.Path.Mount == "" {
		cfg.Path.Mount = "/auth"
	}

	if cfg.Path.TokenKey == "" {
		cfg.Path.TokenKey = "jwt"
	}

	if cfg.Path.HeaderKey == "" {
		cfg.Path.HeaderKey = "X-Token"
	}

	if cfg.Response.DefaultResultType == "" {
		cfg.Response.DefaultResultType = RspTypeJSON
	}

	if cfg.Core.DB == nil {
		return nil, errors.New("缺少Core.DB数据库连接定义！")
	}

	if cfg.Core.Router == nil {
		return nil, errors.New("缺少Core.Router路由定义！")
	}

	if cfg.Core.AuthToken == nil {
		return nil, errors.New("缺少Core.AuthToken定义！")
	}

	ga.Config = cfg
	return ga, nil
}
