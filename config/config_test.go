package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// change workdir to parent to load config.yaml
	err := os.Chdir("..")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name  string
		setup func()
		check func(Config)
	}{
		{
			"defaults",
			func() {
				_ = os.Setenv("U2F_REGISTRATION_TOKEN", "123")
			},
			func(c Config) {
				assert.Equal(t, false, c.Debug)
				assert.Equal(t, "http://localhost:8080", c.URL)
				assert.Equal(t, ":8080", c.Serve)
				assert.Equal(t, &Registration{true, "123"}, c.Registration)
				assert.Equal(t, &Session{"super-secret-key", "localhost", "auth_session"}, c.Session)
				assert.Equal(t, &Db{"storage/database"}, c.Db)
				assert.Equal(t, &Webauthn{"webauthn-demo"}, c.Webauthn)
				assert.NoError(t, c.Validate())
			},
		},
		{
			"validate sqlite path",
			func() {
				_ = os.Setenv("U2F_DB_SQLITE_FILE", "./a/b")
			},
			func(c Config) {
				assert.Equal(t, "./a/b", c.Db.SqliteFile)
				assert.Error(t, c.Validate())
			},
		},
		{
			"validate registration config",
			func() {
			},
			func(c Config) {
				assert.Equal(t, true, c.Registration.Allowed)
				assert.Equal(t, "", c.Registration.Token)
				assert.Error(t, c.Validate())
			},
		},
		{
			"set all configs",
			func() {
				_ = os.Setenv("U2F_DEBUG", "true")
				_ = os.Setenv("U2F_URL", "test-1")
				_ = os.Setenv("U2F_SERVE", "test-2")
				_ = os.Setenv("U2F_REGISTRATION_ALLOWED", "false")
				_ = os.Setenv("U2F_REGISTRATION_TOKEN", "test-3")
				_ = os.Setenv("U2F_SESSION_KEY", "test-4")
				_ = os.Setenv("U2F_SESSION_DOMAIN", "test-5")
				_ = os.Setenv("U2F_SESSION_COOKIE_NAME", "test-5-1")
				_ = os.Setenv("U2F_DB_SQLITE_FILE", "test-6")
				_ = os.Setenv("U2F_WEBAUTHN_RELYING_PARTY_NAME", "test-7")
			},
			func(c Config) {
				assert.Equal(t, true, c.Debug)
				assert.Equal(t, "test-1", c.URL)
				assert.Equal(t, "test-2", c.Serve)
				assert.Equal(t, &Registration{false, "test-3"}, c.Registration)
				assert.Equal(t, &Session{"test-4", "test-5", "test-5-1"}, c.Session)
				assert.Equal(t, &Db{"test-6"}, c.Db)
				assert.Equal(t, &Webauthn{"test-7"}, c.Webauthn)
				assert.NoError(t, c.Validate())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			tt.setup()
			c := InitConfig()
			tt.check(c)
		})
	}
}
