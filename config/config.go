package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// ProjectConfig holds the configuration settings for a single project.
type ProjectConfig struct {
	Name   string   `mapstructure:"name"`
	Port   string   `mapstructure:"port"`
	Proxy  string   `mapstructure:"proxy"`
	Domain []string `mapstructure:"domain"`
}

// AppConfig holds the configuration for the entire application, including multiple projects.
type AppConfig struct {
	Projects []ProjectConfig `mapstructure:"projects"`
}

var once sync.Once

// LoadConfig loads the application configuration from the specified path.
func LoadConfig(path string) (config *AppConfig) {
	once.Do(func() {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			panic("ERROR load config file!")
		}
		log.Println("================ Loaded Configuration ================")
	})
	return
}
