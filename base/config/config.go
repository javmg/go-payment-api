package config

import "os"

const (
	DEFAULT_DB_NAME     = "payment_db"
	DEFAULT_DB_USERNAME = "payment"
	DEFAULT_DB_PASSWORD = "payment"

	DEFAULT_SERVER_HOST = "localhost"
	DEFAULT_SERVER_PORT = "8080"
)

type Config struct {
	DB     *DBConfig
	Server *ServerConfig
}

type DBConfig struct {
	Name     string
	Username string
	Password string
}

type ServerConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {

	dbName := getEnvParamOrDefault("DB_NAME", DEFAULT_DB_NAME)

	dbUsername := getEnvParamOrDefault("DB_USERNAME", DEFAULT_DB_USERNAME)

	dbPassword := getEnvParamOrDefault("DB_PASSWORD", DEFAULT_DB_PASSWORD)

	dbServerHost := getEnvParamOrDefault("SERVER_HOST", DEFAULT_SERVER_HOST)

	dbServerPort := getEnvParamOrDefault("SERVER_PORT", DEFAULT_SERVER_PORT)

	return &Config{

		DB: &DBConfig{
			Name:     dbName,
			Username: dbUsername,
			Password: dbPassword,
		},

		Server: &ServerConfig{
			Host: dbServerHost,
			Port: dbServerPort,
		},
	}
}

func getEnvParamOrDefault(envParamName string, defaultValue string) string {

	value := os.Getenv(envParamName)

	if len(value) == 0 {
		return defaultValue
	}

	return value
}
