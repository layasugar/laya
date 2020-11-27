package mongodb

type Config struct {
	Open        bool   `toml:"open"`
	DSN         string `toml:"dsn"`
	maxIdleConn uint64 `toml:"maxIdleConn"`
	maxOpenConn uint64 `toml:"maxOpenConn"`
}
