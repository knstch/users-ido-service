package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JwtSecret string `envconfig:"JWT_SECRET"`

	JaegerHost  string `envconfig:"JAEGER_HOST" default:"http://localhost:14268/api/traces"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"service-template"`

	PublicHTTPAddr string `envconfig:"PUBLIC_HTTP_ADDR" default:"8080"`

	KafkaAddr string `envconfig:"KAFKA_ADDR" default:"localhost:9092"`

	PlatformURL string `envconfig:"PLATFORM_URL"`

	GoogleAPI GoogleAPI

	DBConfig    DBConfig
	RedisConfig RedisConfig
}

type GoogleAPI struct {
	GoogleAuthHost          string `envconfig:"GOOGLE_AUTH_HOST"`
	GoogleAPIHost           string `envconfig:"GOOGLE_API_HOST"`
	GoogleOAuthClientSecret string `envconfig:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleClientID          string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleRedirectURI       string `envconfig:"GOOGLE_REDIRECT_URI"`
}

type DBConfig struct {
	Host     string `envconfig:"PG_HOST"`
	Port     string `envconfig:"PG_PORT"`
	User     string `envconfig:"PG_USER"`
	Password string `envconfig:"PG_PASSWORD"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     string `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.User)
}

func (cfg *Config) GetRedisDSN() string {
	return fmt.Sprintf("redis://:%s@%s:%s/0", cfg.RedisConfig.Password, cfg.RedisConfig.Host, cfg.RedisConfig.Port)
}

func GetConfig() (*Config, error) {
	config := &Config{}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitENV(dir string) error {
	if err := godotenv.Load(filepath.Join(dir, ".env.local")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("godotenv.Load: %w", err)
		}
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}
	return nil
}
