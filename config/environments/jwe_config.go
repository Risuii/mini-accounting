package environments

import "time"

type JWEConfig struct {
	SecretKey      string        `env-required:"true" env:"JWE_SECRET_KEY"`
	ExpiryDuration time.Duration `env-required:"true" env:"JWE_EXPIRY_DURATION"`
}
