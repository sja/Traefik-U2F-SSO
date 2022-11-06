package config

import "fmt"

type Registration struct {
	Allowed bool   `mapstructure:"allowed"`
	Token   string `mapstructure:"token"`
}

var defaultRegistration = &Registration{
	Allowed: true,
	Token:   "",
}

func (r Registration) Validate() error {
	// Check if registration token is set if registration is allowed
	if r.Allowed && len(r.Token) <= 0 {
		return fmt.Errorf("config error: registration token has to be set, if registration is allowed")
	}
	return nil
}
