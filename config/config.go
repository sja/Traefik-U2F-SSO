package config

import (
	"fmt"
	. "github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
)

type Config struct {
	Debug        bool          `mapstructure:"debug"`
	URL          string        `mapstructure:"url"`
	Serve        string        `mapstructure:"serve"`
	Registration *Registration `mapstructure:"registration"`
	Session      *Session      `mapstructure:"session"`
	Db           *Db           `mapstructure:"db"`
	Webauthn     *Webauthn     `mapstructure:"webauthn"`
}

func InitConfig() Config {
	SetConfigName("config")
	SetConfigType("yaml")
	AddConfigPath(".")
	SetEnvPrefix("U2F")                              // env var prefix, use e.g. U2F_DEBUG
	SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // U2F_DB.SQLITE_FILE -> U2F_DB_SQLITE_FILE
	AutomaticEnv()

	if err := ReadInConfig(); err != nil {
		if _, ok := err.(ConfigFileNotFoundError); !ok {
			log.Fatal(fmt.Errorf("cannot read config: %w", err))
		}
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

func (c *Config) Validate() error {
	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("config error: URL not parsable: %q", c.URL)
	}

	if err := c.Db.Validate(); err != nil {
		return err
	}

	if err := c.Registration.Validate(); err != nil {
		return err
	}

	return nil
}
