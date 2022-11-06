package config

type Webauthn struct {
	RelyingPartyName string `mapstructure:"relying_party_name"`
}

var defaultWebauthn = &Webauthn{RelyingPartyName: "webauthn-demo"}
