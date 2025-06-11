package environments

import "time"

type DbConfigModel struct {
	MaxIdleConns    int           `env-required:"true" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env-required:"true" env:"DB_MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env-required:"true" env:"DB_MAX_LIFE_TIME"`
}
