module telegram

go 1.22.0

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/mentalisit/logger v0.0.0-20240221024243-6f28067f593e
	golang.org/x/image v0.15.0
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => ./telegram-bot-api/
