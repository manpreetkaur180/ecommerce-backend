package config

import "os"

type SMTPConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
}

func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Email:    os.Getenv("SMTP_EMAIL"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}