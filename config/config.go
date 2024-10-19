package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	HttpAuthAddress string `mapstructure:"HTTP_AUTH_SERVER_ADDRESS"`

	Environment string `mapstructure:"NODE_ENV"`
	// db
	PostgresDriver string `mapstructure:"POSTGRES_DRIVER"`
	PostgresSource string
	DbName         string
	DbUser         string
	DbPassword     string
	DbHost         string
	DbPort         string
	DbSSLMode      string
	MigrationURL   string `mapstructure:"MIGRATION_URL"`

	Port string `mapstructure:"PORT"`

	Origin string `mapstructure:"ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	SessionDuration        int           `mapstructure:"SESSION_DURATION"`

	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	HTTPClientAddress string `mapstructure:"HTTP_CLIENT_ADDRESS"`

	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	// SMTP
	SMTPAuthAddress     string `mapstructure:"SMTP_AUTH_ADDRESS"`
	SMTPServerAddress   string `mapstructure:"SMTP_SERVER_ADDRESS"`
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`

	// Redis
	RedisAddress string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	nodeEnv := config.Environment

	// redis
	config.RedisAddress = viper.GetString(fmt.Sprintf("%s_REDIS_ADDRESS", nodeEnv))

	// postgres
	config.DbName = viper.GetString(fmt.Sprintf("%s_POSTGRES_DB", nodeEnv))
	config.DbUser = viper.GetString(fmt.Sprintf("%s_POSTGRES_USER", nodeEnv))
	config.DbPassword = viper.GetString(fmt.Sprintf("%s_POSTGRES_PASSWORD", nodeEnv))
	config.DbHost = viper.GetString(fmt.Sprintf("%s_POSTGRES_HOST", nodeEnv))
	config.DbPort = viper.GetString(fmt.Sprintf("%s_POSTGRES_PORT", nodeEnv))
	config.DbSSLMode = viper.GetString(fmt.Sprintf("%s_SSL_MODE", nodeEnv))

	config.PostgresSource = fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s", "postgres", config.DbUser, config.DbPassword,
		config.DbHost, config.DbPort, config.DbName, config.DbSSLMode)

	return
}
