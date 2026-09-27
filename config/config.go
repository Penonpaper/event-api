package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type JWTConfig struct {
	AccessTokenSecret  string        `mapstructure:"access_token_secret"`
	RefreshTokenSecret string        `mapstructure:"refresh_token_secret"`
	AccessTokenTTL     time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL    time.Duration `mapstructure:"refresh_token_ttl"`
}
type AppConfig struct {
	Port    string
	AppEnv  string
	GinMode string
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"db_name"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host string
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("postgres.password", "DB_PASSWORD")
	viper.BindEnv("jwt.access_token_secret", "ACCESS_TOKEN_SECRET")
	viper.BindEnv("jwt.refresh_token_secret", "REFRESH_TOKEN_SECRET")

	//viper.BindEnv("postgres.db_name", "DB_NAME")
	//viper.BindEnv("jwt.secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла конфигурации: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Ошибка парсинга конфиурации в структуру: %w", err)
	}

	return &cfg, nil

}
