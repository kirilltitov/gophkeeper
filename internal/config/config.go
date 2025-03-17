package config

import (
	"os"
	"strconv"
)

// Config contains all configuration fields of the service.
type Config struct {
	ServerAddress string // Protocol, host and port for HTTP server to bind.
	DatabaseDSN   string // A DSN for database connection.
	TLSCertFile   string // TLS certificate file path.
	TLSKeyFile    string // TLS key file path.
	JWTCookieName string // Auth cookie name.
	JWTSecret     string // A secret for JWT signing.
	JWTTimeToLive int    // Time (in seconds) for JWT expiration configuration.
}

// New creates and returns a new fully set config.
func New() *Config {
	ParseFlags()

	return NewWithoutParsing()
}

// NewWithoutParsing creates and returns a new config without parsing cmd flags.
func NewWithoutParsing() *Config {
	return &Config{
		ServerAddress: getServerAddress(),
		DatabaseDSN:   getDatabaseDSN(),
		TLSCertFile:   getTLSCertFile(),
		TLSKeyFile:    getTLSKeyFile(),
		JWTCookieName: getJWTCookieName(),
		JWTSecret:     getJWTSecret(),
		JWTTimeToLive: getJWTTimeToLive(),
	}
}

// IsTLSEnabled returns true if both TLS cert/key are provided.
func (c *Config) IsTLSEnabled() bool {
	return c.TLSKeyFile != "" && c.TLSCertFile != ""
}

func getServerAddress() string {
	var result = flagBind

	envValue := os.Getenv("RUN_ADDRESS")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getDatabaseDSN() string {
	var result = flagDatabaseDSN

	envValue := os.Getenv("DATABASE_URI")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getTLSCertFile() string {
	var result = tlsCertFile

	envValue := os.Getenv("TLS_CERT_FILE")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getTLSKeyFile() string {
	var result = tlsKeyFile

	envValue := os.Getenv("TLS_KEY_FILE")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getJWTCookieName() string {
	var result = jwtCookieName

	envValue := os.Getenv("JWT_COOKIE_NAME")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getJWTSecret() string {
	var result = jwtSecret

	envValue := os.Getenv("JWT_SECRET")
	if envValue != "" {
		result = envValue
	}

	return result
}

func getJWTTimeToLive() int {
	var result = jwtTimeToLive

	envValue := os.Getenv("JWT_TTL")
	if envValue != "" {
		res, err := strconv.Atoi(envValue)
		if err != nil {
			res = 0
		}
		result = res
	}

	return result
}
