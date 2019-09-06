package auth

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"regexp"
)

type GinAuth struct {
	Config    Config
	Handlers  GinAuthHandlers
	providers map[string]AuthProvider
	renders   map[string]AuthRender
}

func (a *GinAuth) GetDb() *gorm.DB {
	return a.Config.Core.DB
}

func (a *GinAuth) RegisterProvider(p AuthProvider, r AuthRender) error {

	var (
		name = p.GetName()
	)

	// 校验name值是否符合规定
	reg := regexp.MustCompile(`^[0-9a-z_]+$`)
	if !reg.MatchString(name) || len(name) > 16 {
		return fmt.Errorf("Provider[%s]不符合命名规范，请使用小写字母、数字和下划线，长度不超过16个字符！", name)
	}

	if a.HasProvider(name) {
		return fmt.Errorf("Provider[%s]已被注册！", name)
	}

	if regErr := p.OnProviderRegister(a); regErr != nil {
		return regErr
	}

	if regErr := r.OnRenderRegister(a); regErr != nil {
		return regErr
	}

	if a.providers == nil {
		a.providers = map[string]AuthProvider{}
	}

	if a.renders == nil {
		a.renders = map[string]AuthRender{}
	}

	a.providers[name] = p
	a.renders[name] = r

	return nil
}

func (a *GinAuth) HasProvider(name string) bool {
	if a.providers == nil {
		return false
	}
	_, ok := a.providers[name]
	return ok
}

func (a *GinAuth) GetProvider(name string) AuthProvider {
	if a.providers == nil {
		return nil
	}
	return a.providers[name]
}

func (a *GinAuth) GetRender(name string) AuthRender {
	if a.renders == nil {
		return nil
	}
	return a.renders[name]
}
