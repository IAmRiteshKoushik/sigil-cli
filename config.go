package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	RabbitMQURL string `koanf:"rabbitmq_url"`
	SMTPHost    string `koanf:"smtp_host"`
	SMTPPort    int    `koanf:"smtp_port"`
	SMTPUser    string `koanf:"smtp_username"`
	SMTPPass    string `koanf:"smtp_password"`
	SMTPFrom    string `koanf:"smtp_from"`
}

var cfg Config
var k = koanf.New(".")

func LoadConfig() {
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
}

func GetConfig() Config {
	return cfg
}
