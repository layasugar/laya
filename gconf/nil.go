package gconf

const Nil = ConfigError("config: not found")

type ConfigError string

func (e ConfigError) Error() string { return string(e) }
