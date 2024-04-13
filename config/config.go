package config

type Config struct {
	PRIVATE_KEY      string `mapstructure:"PRIVATE_KEY"`
	NODE_URL         string `mapstructure:"NODE_URL"`
	CONTRACT_ADDRESS string `mapstructure:"CONTRACT_ADDRESS"`
	PORT             string `mapstructure:"PORT"`
	ORIGIN           string `mapstructure:"ORIGIN"`
	MODE             string `mapstructure:"MODE"`
	API_HOST         string `mapstructure:"API_HOST"`
	API_VERSION      string `mapstructure:"API_VERSION"`
}
