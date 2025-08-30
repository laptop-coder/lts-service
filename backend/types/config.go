package types

type AppConfig struct {
	DevMode string
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

type SSLConfig struct {
	PathToCert string
	PathToKey  string
}

type StorageConfig struct {
	PathTo string
}

type ValkeyConfig struct {
	Host string
	Port int
}

type Config struct {
	App     AppConfig
	DB      DBConfig
	Env     EnvConfig
	Logs    LogsConfig
	RSA     RSAConfig
	SSL     SSLConfig
	Storage StorageConfig
	Valkey  ValkeyConfig
}
