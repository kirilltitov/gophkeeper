package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerAddress string
	DatabaseDSN   string
	TLSCertFile   string
	TLSKeyFile    string
	JWTCookieName string
	JWTSecret     string
	JWTTimeToLive int
}

func New() *Config {
	ParseFlags()

	return NewWithoutParsing()
}

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
