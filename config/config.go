package RobinConfig

type Configurations struct {
	Name       string
	Technique  string
	EncryptKey string
	Service    ServiceConfig
	Startup    StartupConfig
	Registry   RegistryConfig
}

type ServiceConfig struct {
	Name        string
	Description string
	Path        string
	DLL         bool
	RunAs       string
	Password    string
}

type StartupConfig struct {
	Name string
	Path string
}

type RegistryConfig struct {
	Key        string
	Path       string
	Hive       string
	CustomName string
}
