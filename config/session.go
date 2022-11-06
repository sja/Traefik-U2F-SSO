package config

type Session struct {
	Key        string `mapstructure:"key"`
	Domain     string `mapstructure:"domain"`
	CookieName string `mapstructure:"cookie_name"`
}

var defaultSession = &Session{
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	Key: "super-secret-key",
	// where to bind auth cookie to
	Domain:     "localhost",
	CookieName: "auth_session",
}
