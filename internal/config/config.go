package config

import (
	"github.com/joho/godotenv"
	"os"
)

const ENV_NAME = ".env"

type Config struct{}

func NewConfig() *Config {
	if err := godotenv.Load(ENV_NAME); err != nil {
		panic(err)
	}
	return &Config{}
}

func (e *Config) GetMongoHost() string {
	return os.Getenv("MONGO_HOST")
}

func (e *Config) GetMongoPort() string {
	return os.Getenv("MONGO_PORT")
}

func (e *Config) GetMongoDatabase() string {
	return os.Getenv("MONGO_DATABASE")
}

func (e *Config) GetMongoUser() string {
	return os.Getenv("MONGO_USER")
}

func (e *Config) GetMongoPwd() string {
	return os.Getenv("MONGO_PWD")
}

func (e *Config) GetTgBotName() string {
	return os.Getenv("TG_BOT_NAME")
}

func (e *Config) GetTgBotApiKey() string {
	return os.Getenv("TG_BOT_API_KEY")
}
