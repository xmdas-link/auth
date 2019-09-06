package session

import (
	"github.com/adam-hanna/sessions"
	"github.com/adam-hanna/sessions/auth"
	"github.com/adam-hanna/sessions/store"
	"github.com/adam-hanna/sessions/transport"
)

func New(cfg Config) (*Module, error) {

	var (
		seshStore = store.New(store.Options{
			ConnectionAddress: cfg.RedisAddress,
		})
		seshAuth, err = auth.New(auth.Options{
			Key: []byte(cfg.AuthKey),
		})
		seshTransport = transport.New(transport.Options{
			HTTPOnly: !cfg.Https,
			Secure:   cfg.Secure, // note: can't use secure cookies in development!
		})

		m = &Module{}
	)

	if err != nil {
		return nil, err
	}

	m.Service = sessions.New(seshStore, seshAuth, seshTransport, sessions.Options{})

	return m, nil
}
