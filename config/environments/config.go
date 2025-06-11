package environments

type ConfigModel struct {
	App          AppConfigModel
	JWE          JWEConfig
	DB           DbConfigModel
	DBAccounting DBAccounting
}
