package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	ConfPort                = "port"
	ConfRegistrationAllowed = "registrationAllowed"
	ConfRegistrationToken   = "registrationToken"
	ConfURL                 = "url"
	ConfDebug               = "debug"
	ConfDomain              = "domain"
	ConfSqliteFile          = "sqliteFile"
	ConfSessionKey          = "sessionKey"
	ConfRelyingPartyName    = "webauthRelyingPartyName"
)

func setDefaultVlaues() {
	viper.SetDefault(ConfPort, 8080)
	viper.SetDefault(ConfRegistrationAllowed, false)
	viper.SetDefault(ConfURL, "https://localhost:8080")
	viper.SetDefault(ConfDebug, false)
	viper.SetDefault(ConfDomain, "localhost")
	viper.SetDefault(ConfSqliteFile, "storage/database")
	viper.SetDefault(ConfSessionKey, "super-secret-key") // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	viper.SetDefault(ConfRelyingPartyName, "webauthn-demo")
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("u2f") // env var prefix, use e.g. U2F_PORT
	setDefaultVlaues()
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	return nil
}
