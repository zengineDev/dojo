package dojo

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Port        int    `json:"port"`
	Environment string `json:"environment"`
	Domain      string `json:"domain"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ViewConfig struct {
	Path string `json:"path"`
}

type AssetsConfigs struct {
	Path string `json:"path"`
}

type SessionConfig struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type AuthenticationProvider string

const (
	OAuthAuthenticationProvider    AuthenticationProvider = "oauth"
	DatabaseAuthenticationProvider AuthenticationProvider = "database"
)

type AuthenticationConfig struct {
	Provider AuthenticationProvider
	// The Configuration for the database provider
	Table string `json:"table"`
	// The Configuration for the oauth provider
	Endpoint     string   `json:"endpoint"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Scopes       []string `json:"scopes"`
	RedirectPath string   `json:"redirectPath"`
}

type DefaultConfiguration struct {
	App     AppConfig            `yaml:"app"`
	DB      DatabaseConfig       `yaml:"db"`
	View    ViewConfig           `yaml:"view"`
	Assets  AssetsConfigs        `yaml:"assets"`
	Session SessionConfig        `yaml:"session"`
	Auth    AuthenticationConfig `yaml:"auth"`
}

func LoadConfigs(cfg *DefaultConfiguration) *DefaultConfiguration {
	v := viper.New()

	// Viper settings
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("$CONFIG_DIR/")

	// Environment variable settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	// Global configuration
	v.SetDefault("shutdownTimeout", 15*time.Second)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		v.SetDefault("no_color", true)
	}

	err := v.ReadInConfig()
	if err != nil {
		panic(errors.Wrap(err, "Cant read configuration file"))
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(errors.Wrap(err, "Cant unmarshall configuration"))
	}

	return cfg
}
