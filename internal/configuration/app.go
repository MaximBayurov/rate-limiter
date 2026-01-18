package configuration

type AppConf struct {
	LoginAttempts    int `mapstructure:"login_attempts"`
	PasswordAttempts int `mapstructure:"password_attempts"`
	IPAttempts       int `mapstructure:"ip_attempts"`
}
