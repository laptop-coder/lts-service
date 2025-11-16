package types

type AppConfig struct {
	PortBackend  string
	PortFrontend string
	BcryptCost int
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

type StorageConfig struct {
	PathTo string
}

type Config struct {
	App     AppConfig
	DB      DBConfig
	Env     EnvConfig
	Logs    LogsConfig
	RSA     RSAConfig
	Storage StorageConfig
}
