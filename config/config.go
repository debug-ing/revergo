package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type ProjectConfig struct {
	Name  string `mapstructure:"name"`
	Port  string `mapstructure:"port"`
	Proxy string `mapstructure:"proxy"`
}

type AppConfig struct {
	Projects []ProjectConfig `mapstructure:"projects"`
}

var once sync.Once

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
