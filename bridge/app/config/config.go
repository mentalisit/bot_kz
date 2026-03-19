package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigBot struct {
	Logger struct {
		Token   string `yaml:"token"`
		ChatId  int64  `yaml:"chat_id"`
		Webhook string `yaml:"webhook"`
	} `yaml:"logger"`
	Postgress struct {
		Host     string `yaml:"host" env-default:"127.0.0.1:3306"`
		Name     string `yaml:"name" env-default:"rsbot"`
		Username string `yaml:"username" env-default:"root"`
		Password string `yaml:"password" env-default:"root"`
	} `yaml:"postgress"`
	Matrix struct {
		HomeserverURL string `yaml:"homeserverURL" env-default:"http://10.0.0.184:6167"`
		Username      string `yaml:"username" env-default:"@bridge_bot:matrix.mentalisit.myds.me"`
		Password      string `yaml:"password" env-default:"botPassword"`
		ASToken       string `yaml:"asToken"`
		HSToken       string `yaml:"hsToken"`
	} `yaml:"matrix"`
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
