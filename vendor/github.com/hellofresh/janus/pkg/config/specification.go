package config

import (
	"time"

	"github.com/hellofresh/logging-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT"`
	Debug                bool          `envconfig:"DEBUG"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL"`
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD"`
	Log                  logging.LogConfig
	Web                  Web
	Database             Database
	Storage              Storage
	Stats                Stats
	Tracing              Tracing
	TLS                  TLS
}

// Web represents the API configurations
type Web struct {
	Port        int  `envconfig:"API_PORT"`
	ReadOnly    bool `envconfig:"API_READONLY"`
	Credentials Credentials
	TLS         TLS
}

// TLS represents the TLS configurations
type TLS struct {
	Port     int    `envconfig:"PORT"`
	CertFile string `envconfig:"CERT_PATH"`
	KeyFile  string `envconfig:"KEY_PATH"`
	Redirect bool   `envconfig:"REDIRECT"`
}

// IsHTTPS checks if you have https enabled
func (s *TLS) IsHTTPS() bool {
	return s.CertFile != "" && s.KeyFile != ""
}

// Storage holds the configuration for a storage
type Storage struct {
	DSN string `envconfig:"STORAGE_DSN"`
}

// Database holds the configuration for a database
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
}

// Stats holds the configuration for stats
type Stats struct {
	DSN                   string   `envconfig:"STATS_DSN"`
	Prefix                string   `envconfig:"STATS_PREFIX"`
	IDs                   string   `envconfig:"STATS_IDS"`
	AutoDiscoverThreshold uint     `envconfig:"STATS_AUTO_DISCOVER_THRESHOLD"`
	AutoDiscoverWhiteList []string `envconfig:"STATS_AUTO_DISCOVER_WHITE_LIST"`
	ErrorsSection         string   `envconfig:"STATS_ERRORS_SECTION"`
}

// Credentials represents the credentials that are going to be
// used by admin JWT configuration
type Credentials struct {
	// Algorithm defines admin JWT signing algorithm.
	// Currently the following algorithms are supported: HS256, HS384, HS512.
	Algorithm string `envconfig:"ALGORITHM"`
	Secret    string `envconfig:"SECRET"`
	Github    Github
	Basic     Basic
}

// Basic holds the basic users configurations
type Basic struct {
	Users map[string]string `envconfig:"BASIC_USERS"`
}

// Github holds the github configurations
type Github struct {
	Organizations []string          `envconfig:"GITHUB_ORGANIZATIONS"`
	Teams         map[string]string `envconfig:"GITHUB_TEAMS"`
}

// IsConfigured checks if github is enabled
func (auth *Github) IsConfigured() bool {
	return len(auth.Organizations) > 0 ||
		len(auth.Teams) > 0
}

// GoogleCloudTracing holds the Google Application Default Credentials
type GoogleCloudTracing struct {
	ProjectID    string `envconfig:"TRACING_GC_PROJECT_ID"`
	Email        string `envconfig:"TRACING_GC_EMAIL"`
	PrivateKey   string `envconfig:"TRACING_GC_PRIVATE_KEY"`
	PrivateKeyID string `envconfig:"TRACING_GC_PRIVATE_ID"`
}

// AppdashTracing holds the Appdash tracing configuration
type AppdashTracing struct {
	DSN string `envconfig:"TRACING_APPDASH_DSN"`
	URL string `envconfig:"TRACING_APPDASH_URL"`
}

// Tracing represents the distributed tracing configuration
type Tracing struct {
	GoogleCloudTracing GoogleCloudTracing `mapstructure:"googleCloud"`
	AppdashTracing     AppdashTracing     `mapstructure:"appdash"`
}

// IsGoogleCloudEnabled checks if google cloud is enabled
func (t Tracing) IsGoogleCloudEnabled() bool {
	return len(t.GoogleCloudTracing.Email) > 0 && len(t.GoogleCloudTracing.PrivateKey) > 0 && len(t.GoogleCloudTracing.PrivateKeyID) > 0 && len(t.GoogleCloudTracing.ProjectID) > 0
}

// IsAppdashEnabled checks if appdash is enabled
func (t Tracing) IsAppdashEnabled() bool {
	return len(t.AppdashTracing.DSN) > 0
}

func init() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("tls.port", "8433")
	viper.SetDefault("tls.redirect", true)
	viper.SetDefault("backendFlushInterval", "20ms")
	viper.SetDefault("database.dsn", "file:///etc/janus")
	viper.SetDefault("storage.dsn", "memory://localhost")
	viper.SetDefault("web.port", "8081")
	viper.SetDefault("web.tls.port", "8444")
	viper.SetDefault("web.tls.redisrect", true)
	viper.SetDefault("web.credentials.algorithm", "HS256")
	viper.SetDefault("web.credentials.basic.users", map[string]string{
		"admin": "admin",
	})
	viper.SetDefault("stats.dsn", "log://")
	viper.SetDefault("stats.errorsSection", "error-log")

	logging.InitDefaults(viper.GetViper(), "log")
}

//Load configuration variables
func Load(configFile string) (*Specification, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("janus")
		viper.AddConfigPath("/etc/janus")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "No config file found")
	}

	var config Specification
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

//LoadEnv loads configuration from environment variables
func LoadEnv() (*Specification, error) {
	var config Specification

	// ensure the defaults are loaded
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
