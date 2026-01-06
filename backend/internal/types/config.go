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

type LogsConfig struct {
	PathToBackend string
}

type PostgresConfig struct {
	DBName   string
	Password string
	Port     string
	SSLMode  string // "disable" or "enable"
	User     string
}

type RSAConfig struct {
	PathToPrivateKey string
	PathToPublicKey  string
}

type RolesConfig struct {
	User      string
	Moderator string
}

type StorageConfig struct {
	PathTo string
}

type Config struct {
	App      AppConfig
	Bcrypt   BcryptConfig
	Env      EnvConfig
	Logs     LogsConfig
	Postgres PostgresConfig
	RSA      RSAConfig
	Role     RolesConfig
	Storage  StorageConfig
}
