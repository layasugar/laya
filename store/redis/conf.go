package redis

type Config struct {
	Open        bool   `toml:"open"`
	DB          int    `toml:"db"`
	Addr        string `toml:"addr"`
	Pwd         string `toml:"pwd"`
	PoolSize    int    `toml:"poolSize"`
	MaxRetries  int    `toml:"maxRetries"`
	IdleTimeout int    `toml:"idleTimeout"`
}
