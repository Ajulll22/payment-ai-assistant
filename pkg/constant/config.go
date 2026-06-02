package constant

import (
	"os"
	"strconv"
)

type Config struct {
	App               AppConfig
	AssistantSQL      AssistantConfig
	AssistantAnalytic AssistantConfig
	DB                DBConfig
}

func GetEnv() *Config {
	dbTimeout, _ := strconv.Atoi(os.Getenv("DB_TIMEOUT"))

	return &Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: os.Getenv("APP_PORT"),
			Key:  os.Getenv("APP_KEY"),
		},
		AssistantSQL: AssistantConfig{
			Model: os.Getenv("ASSISTANT_SQL_MODEL"),
			URL:   os.Getenv("ASSISTANT_SQL_URL"),
		},
		AssistantAnalytic: AssistantConfig{
			Model: os.Getenv("ASSISTANT_ANALYTIC_MODEL"),
			URL:   os.Getenv("ASSISTANT_ANALYTIC_URL"),
		},
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Timeout:  dbTimeout,
		},
	}
}

type AppConfig struct {
	Name string
	Port string
	Key  string
}

type AssistantConfig struct {
	Model string
	URL   string
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	Timeout  int
}
