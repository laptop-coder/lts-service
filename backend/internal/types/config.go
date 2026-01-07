// Package types provides types to use in app (not in db, e.g.) In this package
// there are types that are not models.
package types

type AppConfig struct {
	PortBackend  string
	PortFrontend string
}

type BcryptConfig struct {
	Cost int
}

type EnvConfig struct {
	PathTo string
}

type PostgresConfig struct {
	DBName   string
	Host     string
	Password string
	Port     string
	SSLMode  string // "disable" or "enable"
	TimeZone string
	User     string
}

type RSAConfig struct {
	PathToPrivateKey string
	PathToPublicKey  string
}

type StorageConfig struct {
	PathTo string
}
