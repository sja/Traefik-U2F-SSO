package internal

import (
	"github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	_ "github.com/koesie10/webauthn/attestation"
	"github.com/koesie10/webauthn/webauthn"
	"github.com/spf13/viper"
)

func InitWebauthn(storage *storage.Storage) (*webauthn.WebAuthn, error) {
	return webauthn.New(&webauthn.Config{
		RelyingPartyName:   viper.GetString(config.ConfRelyingPartyName),
		AuthenticatorStore: storage,
		Debug:              viper.GetBool(config.ConfDebug),
	})
}
