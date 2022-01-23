package util

import (
	"time"

	"github.com/spf13/viper"
)

// Application configuration
type Config struct {
	DBDriver              string        `mapstructure:"DB_DRIVER"`
	DBSource              string        `mapstructure:"DB_SOURCE"`
	ServerAddress         string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file and env vars
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// Read env vars if they exist (overrides matching vars from the file)
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
