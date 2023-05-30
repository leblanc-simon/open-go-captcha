package config

type Config struct {
	Captcha struct {
		SecretKey	string	`yaml:"secret_key" env:"CAPTCHA_SECRET_KEY" env-default:"" env-description:"Secret key to crypt user tokens. Required !"`
		Min			int		`yaml:"min" env:"CAPTCHA_MIN" env-default:"5" env-description:"Minimum number of icons display in captcha"`
		Max			int		`yaml:"max" env:"CAPTCHA_MAX" env-default:"8" env-description:"Maximum number of icons display in captcha"`
		MaxGood		int		`yaml:"max_good" env:"CAPTCHA_MAX_GOOD" env-default:"2" env-description:"Maximum number of icons for the good answer"`
		Rotate		bool	`yaml:"rotate" env:"CAPTCHA_ROTATE" env-default:"1" env-description:"Rotate icons in captcha"`
		Flip		bool	`yaml:"flip" env:"CAPTCHA_FLIP" env-default:"1" env-description:"Flip icons in captcha"`
	} `yaml:"captcha"`

	Redis struct {
		Host		string	`yaml:"host" env:"REDIS_HOST" env-default:"localhost" env-description:"Host for Redis database"`
		Port		int		`yaml:"port" env:"REDIS_PORT" env-default:"6379" env-description:"Port for Redis database"`
		Password	string	`yaml:"password" env:"REDIS_PASSWORD" env-default:"" env-description:"Password for Redis database"`
		Db			int		`yaml:"db" env:"REDIS_DB" env-default:"0" env-description:"Database for Redis database"`
		Expire		int		`yaml:"expire" env:"REDIS_EXPIRE" env-default:"300" env-description:"Time before captcha result expire"` // 5 minutes
		LongExpire	int		`yaml:"long_expire" env:"REDIS_LONG_EXPIRE" env-default:"1800" env-description:"Time before captcha validation expire"` // 30 minutes
		KeyPrefix	string	`yaml:"key_prefix" env:"REDIS_KEY_PREFIX" env-default:"opc_" env-description:"Prefix to the key store in Redis"`
	} `yaml:"redis"`

	Server struct {
		Host		string	`yaml:"host" env:"SERVER_HOST" env-default:"127.0.0.1" env-description:"Listen IP for web server"`
		Port		int		`yaml:"port" env:"SERVER_PORT" env-default:"4242" env-description:"Listen port for web server"`
		LogLevel	string	`yaml:"log_level" env:"LOG_LEVEL" env-default:"error" env-description:"Log level"`
	} `yaml:"server"`
}
