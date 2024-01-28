package config

import (
	"os"
	"schedule/logs"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis RedisConfig            `yaml:"redis"`
	MySQL map[string]MySQLConfig `yaml:"mysql"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MySQLConfig struct {
	Name     string `yaml:"name"`
	PassWord string `yaml:"password"`
	Addr     string `yaml:"addr"`
	DB       string `yaml:"db"`
}

func InitServerConfig(fileName string) {
	configFile, err := os.ReadFile(fileName)
	if err != nil {
		logs.GetInstance().Logger.Errorf("server config %s not found %s", fileName, err)
		return
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		logs.GetInstance().Logger.Errorf("config yaml unmarshal error %s", err)
	}
}

var config = &Config{}

func GetServerConfig() *Config {
	return config
}
