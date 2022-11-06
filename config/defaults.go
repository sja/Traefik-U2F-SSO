package config

// Default config, can be overridden by env variables
var defaultConfig = &Config{
	Debug: false,

	// own url, used for redirects
	URL:   "http://localhost:8080",
	Serve: ":8080",

	Registration: defaultRegistration,
	Session:      defaultSession,
	Db:           defaultDB,
	Webauthn:     defaultWebauthn,
}
