package pkg

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

const CONFIG_PATH = "configs/env.toml"

var cfg Config
var k = koanf.New(".")

type Config struct {
	RabbitMQURL string `koanf:"rabbitmq_url"`
	StorageDir  string `koanf:"storage_dir"`
	SMTPHost    string `koanf:"smtp_host"`
	SMTPPort    int    `koanf:"smtp_port"`
	SMTPUser    string `koanf:"smtp_username"`
	SMTPPass    string `koanf:"smtp_password"`
}

func LoadConfig() {
	if err := k.Load(file.Provider(CONFIG_PATH), toml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
}

func GetConfig() Config {
	return cfg
}
