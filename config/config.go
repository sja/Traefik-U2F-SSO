package config

import (
	"fmt"
	. "github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Registration struct {
	Allowed bool   `mapstructure:"allowed"`
	Token   string `mapstructure:"token"`
}

type Session struct {
	Key    string `mapstructure:"key"`
	Domain string `mapstructure:"domain"`
}

type Db struct {
	SqliteFile string `mapstructure:"sqlite_file"`
}

type Webauthn struct {
	RelyingPartyName string `mapstructure:"relying_party_name"`
}

type Config struct {
	Debug        bool         `mapstructure:"debug"`
	URL          string       `mapstructure:"url"`
	Serve        string       `mapstructure:"serve"`
	Registration Registration `mapstructure:"registration"`
	Session      Session      `mapstructure:"session"`
	Db           Db           `mapstructure:"db"`
	Webauthn     Webauthn     `mapstructure:"webauthn"`
}

func (c *Config) Validate() error {
	// Check if registration token is set if registration is allowed
	if c.Registration.Allowed && len(c.Registration.Token) <= 0 {
		return fmt.Errorf("config error: registration token has to be set, if registration is allowed")
	}

	// Check if the sqlite db path exists, db files are created automatically
	dir, _ := filepath.Abs(filepath.Dir(c.Db.SqliteFile))
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("config error: path to db does not exist: %v", c.Db.SqliteFile)
	}
	return nil
}

func InitConfig() Config {
	SetConfigName("config")
	SetConfigType("yaml")
	AddConfigPath(".")
	SetEnvPrefix("U2F")                              // env var prefix, use e.g. U2F_DEBUG
	SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // U2F_DB.SQLITE_FILE -> U2F_DB_SQLITE_FILE
	AutomaticEnv()

	if err := ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("cannot read config:: %w", err))
	}

	var config Config
	if err := Unmarshal(&config); err != nil {
		log.Fatal(fmt.Errorf("cannot decode config: %w", err))
	}

	if config.Debug {
		log.Printf("Loaded config file %q and ENVs U2F_*, result:\n%v",
			ConfigFileUsed(), AllSettings())
	}
	return config
}
