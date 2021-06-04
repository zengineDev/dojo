package dojo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
	Testing     Environment = "testing"
	Staging     Environment = "staging"
)

type AppConfig struct {
	Name        string      `json:"name" yaml:"name"`
	Version     string      `json:"version" yaml:"version"`
	Port        int         `json:"port" yaml:"port"`
	Environment Environment `json:"environment" yaml:"environment"`
	Domain      string      `json:"domain" yaml:"domain"`
	Debug       bool        `json:"debug" yaml:"debug"`
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"`
}

func (c DatabaseConfig) DSN() string {
	maxConLife := time.Hour
	maxConIdle := time.Minute * 30
	health := time.Minute
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s "+
		"pool_max_conns=10 pool_min_conns=1 pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v pool_health_check_period=%v",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode, maxConLife, maxConIdle, health)
}

type ViewConfig struct {
	Path string `json:"path" yaml:"path"`
}

type AssetsConfigs struct {
	Path string `json:"path" yaml:"path"`
}

type SessionConfig struct {
	Name   string `json:"name" yaml:"name"`
	Secret string `json:"secret" yaml:"secret"`
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
	Endpoint     string   `json:"endpoint" yaml:"endpoint"`
	ClientID     string   `json:"clientId" yaml:"client_id"`
	ClientSecret string   `json:"clientSecret" yaml:"client_secret"`
	Scopes       []string `json:"scopes" yaml:"scopes"`
	RedirectPath string   `json:"redirectPath" yaml:"redirect_path"`
}

type RedisConfig struct {
	Host     string `json:"host" yaml:"host"`
	Password string `json:"password" yaml:"password"`
	Database int    `json:"database" yaml:"database"`
	Port     int    `json:"port" yaml:"port"`
}

type DefaultConfiguration struct {
	App     AppConfig            `json:"dojo" yaml:"dojo"`
	DB      DatabaseConfig       `json:"db" yaml:"db"`
	View    ViewConfig           `json:"view" yaml:"view"`
	Assets  AssetsConfigs        `json:"assets" yaml:"assets"`
	Session SessionConfig        `json:"session" yaml:"session"`
	Auth    AuthenticationConfig `json:"auth" yaml:"auth"`
	Redis   RedisConfig          `json:"redis" yaml:"redis"`
}

const defaultShutdownTimeoutSeconds = 15
const defaultPostgresPort = 5432

var once sync.Once

var (
	instance *DefaultConfiguration
)

func LoadConfigs(cfg *DefaultConfiguration) *DefaultConfiguration {
	once.Do(func() {
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
		v.SetDefault("shutdownTimeout", defaultShutdownTimeoutSeconds)
		if _, ok := os.LookupEnv("NO_COLOR"); ok {
			v.SetDefault("no_color", true)
		}

		// Database configuration
		_ = v.BindEnv("db.host")
		v.SetDefault("db.port", defaultPostgresPort)
		_ = v.BindEnv("db.user")
		_ = v.BindEnv("db.password")
		_ = v.BindEnv("db.database")

		err := v.ReadInConfig()
		if err != nil {
			panic(errors.Wrap(err, "Cant read configuration file"))
		}

		err = v.Unmarshal(&cfg)
		if err != nil {
			panic(errors.Wrap(err, "Cant unmarshall configuration"))
		}

		instance = cfg
	})

	return instance

}

func GetConfig() *DefaultConfiguration {
	if instance == nil {
		LoadConfigs(&DefaultConfiguration{})
	}

	return instance
}
