package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `yaml:"Host"`
		Port int    `yaml:"Port"`
	} `yaml:"Server"`

	MySQL struct {
		Host     string `yaml:"Host"`
		Port     int    `yaml:"Port"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		Database string `yaml:"Database"`
	} `yaml:"MySQL"`

	OpenAI struct {
		ApiKey string `yaml:"ApiKey"`
	} `yaml:"OpenAI"`
}

var Conf *Config

func LoadConfig() (*Config, error) {
	if Conf != nil {
		return Conf, nil
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	Conf = &Config{}
	err = viper.Unmarshal(Conf)
	if err != nil {
		return nil, err
	}
	return Conf, nil
}
