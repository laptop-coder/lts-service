package types

type AppConfig struct {
	PortBackend  string
	PortFrontend string
}

type BcryptConfig struct {
	Cost int
}

type DBConfig struct {
	PathTo string
}

type EnvConfig struct {
	PathTo string
}

type LogsConfig struct {
	PathToBackend string
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
	App     AppConfig
	Bcrypt  BcryptConfig
	DB      DBConfig
	Env     EnvConfig
	Logs    LogsConfig
	RSA     RSAConfig
	Role    RolesConfig
	Storage StorageConfig
}
