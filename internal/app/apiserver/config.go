package apiserver

// Config ...
type Config struct {
	BindAddr         string `toml:"bind_addr"`
	LogLevel         string `toml:"log_level"`
	MongoURI         string `toml:"mongo_uri"`
	APIKey           string `toml:"etherscan_api_key"`
	APIURL           string `toml:"etherscan_api_url"`
	DaemonMode       bool   `toml:"daemon_mode"`
	LoadTransactions int    `toml:"load_transactions"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
