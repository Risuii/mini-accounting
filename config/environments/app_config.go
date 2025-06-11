package environments

import "time"

type AppConfigModel struct {
	Port             string        `env-required:"true" env:"APP_PORT"`
	Name             string        `env-required:"true" env:"APP_NAME"`
	SecretKey        string        `env-required:"true" env:"APP_SECRET_KEY"`
	LocationTimezone string        `env-required:"true" env:"APP_LOCATION_TIMEZONE"`
	LoggingTimeout   time.Duration `env-required:"true" env:"APP_LOGGING_TIMEOUT"`
}
