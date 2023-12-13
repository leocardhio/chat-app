package config

var (
	Cfg Config
)

type Config struct {
	RedisHost string `envconfig:"REDIS_HOST" default:"localhost:6379"`
	RedisPwd  string `envconfig:"REDIS_PWD" default:""`

	Port string `envconfig:"HOST_PORT" default:"8080"`
}

func init() {
	Cfg = Config{
		RedisHost: ":6379",
		RedisPwd:  "",
		Port:      "8080",
	}
	// envconfig.MustProcess("", Cfg)
}