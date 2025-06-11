package environments

type DBAccounting struct {
	HOST     string `env-required:"true" env:"DB_Accounting_HOST"`
	PORT     string `env-required:"true" env:"DB_Accounting_PORT"`
	USER     string `env-required:"true" env:"DB_Accounting_USER"`
	PASSWORD string `env-required:"true" env:"DB_Accounting_PASSWORD"`
	DATABASE string `env-required:"true" env:"DB_Accounting_DATABASE"`
}
