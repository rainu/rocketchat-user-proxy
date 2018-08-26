package config

import (
	"github.com/alexflint/go-arg"
)

type Config struct {
	chatConfig

	BindPort int `arg:"--bind-port,env:BIND_PORT,help:The HTTP server bind port."`
}

type chatConfig struct {
	Url           string `arg:"--chat-url,env:CHAT_URL,help:The url of the rocket.chat instance."`
	Username      string `arg:"--chat-user,env:CHAT_USERNAME,help:The username for the chat."`
	PasswordPlain string `arg:"--chat-password,env:CHAT_PASSWORD,help:The user password (plain-text)."`
	PasswordHash  string `arg:"--chat-password-hash,env:CHAT_PASSWORD_HASH,help:The user hashed password (sha-256)."`
}

func New() *Config {
	cfg := &Config{
		chatConfig: chatConfig{},
		BindPort:   8080,
	}

	arg.Parse(cfg)

	return cfg
}
