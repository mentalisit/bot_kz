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
	Port string `yaml:"port"`
}

var Instance *ConfigBot
var once sync.Once

func InitConfig() *ConfigBot {
	once.Do(func() {
		Instance = &ConfigBot{}
		err := cleanenv.ReadConfig("docker/config/config.yml", Instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(Instance, nil)
			fmt.Println(help)
		}
	})
	return Instance
}
