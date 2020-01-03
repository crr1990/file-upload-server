package config

type Config struct {
	Path string `yaml:"path"`
}

var (
	cfg *Config
)

func SetConfig(config *Config) {
	cfg = config
}
func GetConfig() *Config {
	return cfg
}
