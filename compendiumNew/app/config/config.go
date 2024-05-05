package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type ConfigBot struct {
	Logger struct {
		Token   string `yaml:"token"`
		ChatId  int64  `yaml:"chat_id"`
		Webhook string `yaml:"webhook"`
	} `yaml:"logger"`
	Mongo     string `yaml:"mongo"`
	Port      string `yaml:"port"`
	Postgress struct {
		Host     string `yaml:"host" env-default:"127.0.0.1:5432"`
		Name     string `yaml:"name" env-default:"postgres"`
		Username string `yaml:"username" env-default:"postgres"`
		Password string `yaml:"password" env-default:"root"`
	} `yaml:"postgress"`
}

var Instance *ConfigBot
var once sync.Once

func InitConfig() *ConfigBot {
	once.Do(func() {
		Instance = &ConfigBot{}
		err := cleanenv.ReadConfig("config.yml", Instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(Instance, nil)
			fmt.Println(help)
		}
	})
	return Instance
}
