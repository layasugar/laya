package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func load() {
	var config struct{ Name string }
	if _, err := toml.DecodeFile("example.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
}
