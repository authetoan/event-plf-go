package configs

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	App struct {
		Port string `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `mapstructure:"app"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`

	Redis struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"redis"`

	Kafka struct {
		Broker string `mapstructure:"broker"`
	} `mapstructure:"kafka"`
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig(configPath string) *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configPath)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		instance = &Config{}
		if err := viper.Unmarshal(instance); err != nil {
			log.Fatalf("Error unmarshalling config: %v", err)
		}
	})
	return instance
}
