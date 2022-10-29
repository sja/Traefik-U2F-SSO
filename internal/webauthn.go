package internal

import (
	"github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	_ "github.com/koesie10/webauthn/attestation"
	"github.com/koesie10/webauthn/webauthn"
)

func InitWebauthn(c config.Config, storage *storage.Storage) (*webauthn.WebAuthn, error) {
	return webauthn.New(&webauthn.Config{
		RelyingPartyName:   c.Webauthn.RelyingPartyName,
		AuthenticatorStore: storage,
		Debug:              c.Debug,
	})
}
