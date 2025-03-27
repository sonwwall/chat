package configs

type Config struct {
	DatabaseConfig DatabaseConfig
	ZapConfig      ZapConfig
}

type ZapConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type DatabaseConfig struct {
	MysqlConfig MysqlConfig
	RedisConfig RedisConfig
}

type MysqlConfig struct {
	Addr     string
	Username string
	Password string
	DB       string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
