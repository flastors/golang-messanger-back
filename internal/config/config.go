package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "config/config.yaml"

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	ChatService ChatServiceConfig `yaml:"chat_service"`
}

type ServerConfig struct {
	RestAPI RestAPIConfig `yaml:"rest_api"`
}

type RestAPIConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type ChatServiceConfig struct {
	MaxMessageLimit int            `yaml:"max_message_limit"`
	MaxMessageSize  int            `yaml:"max_message_size"`
	MaxTitleSize    int            `yaml:"max_title_size"`
	MinTitleSize    int            `yaml:"min_title_size"`
	ChatDB          DatabaseConfig `yaml:"chat_db"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"CHAT_SERVICE_DB_HOST"`
	Port     int    `yaml:"port" env:"CHAT_SERVICE_DB_PORT"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
	Timezone string `yaml:"timezone"`
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		// Здесь можно инициализировать конфигурацию из файла или переменных окружения
		instance = &Config{}
		readErr := cleanenv.ReadConfig(configPath, instance)
		if readErr != nil {
			description, descrErr := cleanenv.GetDescription(instance, nil)
			if descrErr != nil {
				panic(descrErr)
			}
			slog.Info(description)
			slog.Error(
				"failed to read config",
				slog.String("err", readErr.Error()),
				slog.String("path", configPath),
			)
			os.Exit(1)
		}
	})
	return instance
}
