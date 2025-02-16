package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8081

var flagBind = fmt.Sprintf(":%d", defaultPort)
var flagDatabaseDSN = "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres"
var tlsCertFile = ""
var tlsKeyFile = ""
var jwtCookieName = "access_token"
var jwtSecret = "hesoyam"
var jwtTimeToLive int = 86400

func ParseFlags() {
	flag.StringVar(&flagBind, "bind", flagBind, "Host and port to bind")
	flag.StringVar(&flagDatabaseDSN, "dsn", flagDatabaseDSN, "Database DSN")
	flag.StringVar(&tlsCertFile, "tls_crt", tlsCertFile, "TLS cert file path")
	flag.StringVar(&tlsKeyFile, "tls_key", tlsKeyFile, "TLS key file path")
	flag.StringVar(&jwtCookieName, "jwt_cookie", jwtCookieName, "JWT Cookie name")
	flag.StringVar(&jwtSecret, "jwt_secret", jwtSecret, "JWT Secret")
	flag.IntVar(&jwtTimeToLive, "jwt_ttl", jwtTimeToLive, "JWT Time To Live")

	flag.Parse()
}
