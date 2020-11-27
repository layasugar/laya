package db

type Config struct {
	Driver          string `toml:"driver"`
	Open            bool   `toml:"open"`
	Dsn             string `toml:"dsn"`
	MaxIdleConn     int    `toml:"maxIdleConn"`
	MaxOpenConn     int    `toml:"maxOpenConn"`
	ConnMaxLifetime int    `toml:"connMaxLifetime"`
	ConnTimeout     int    `toml:"connTimeout"`
}
