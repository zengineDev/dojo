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
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Port        int         `json:"port"`
	Environment Environment `json:"environment"`
	Domain      string      `json:"domain"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (c DatabaseConfig) DSN() string {
	maxConLife, _ := time.ParseDuration("10 minutes")
	maxConIdle, _ := time.ParseDuration("15 minutes")
	health, _ := time.ParseDuration("1 minute")
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require"+
		"pool_max_conns=10 pool_min_conns=0 pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v pool_health_check_period=%v",
		c.Host, c.Port, c.User, c.Password, c.Database, maxConLife, maxConIdle, health)
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
