package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type ConfigBot struct {
	Token struct {
		TokenDiscord   string `yaml:"token_discord"`
		TokenTelegram  string `yaml:"token_telegram"`
		NameDbWhatsapp string `yaml:"name_db_whatsapp"`
	} `yaml:"token"`
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
	MySQl   string `yaml:"mysql"`
	WsToken string `yaml:"ws_token"`
}

var Instance *ConfigBot
var once sync.Once

func InitConfig(nameApp string) *ConfigBot {
	once.Do(func() {
		Instance = &ConfigBot{}
		err := cleanenv.ReadConfig("docker/config/config.yml", Instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(Instance, nil)
			fmt.Println(help)
		}
	})
	if Instance != nil {
		l := Instance.Logger
		InitLogger(nameApp, l.Webhook, l.Token, l.ChatId)
	}
	return Instance
}
