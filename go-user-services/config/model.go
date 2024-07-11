package config

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DBIndex  int
}
