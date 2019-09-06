package password

import "golang.org/x/crypto/bcrypt"

func New() *Provider {
	var (
		p = &Provider{
			CryptCost: bcrypt.DefaultCost,
			Name:      Name,
		}
	)

	return p
}
