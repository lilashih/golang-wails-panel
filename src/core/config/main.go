package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/golobby/dotenv"
)

var App AppConfig
var Project ProjectConfig

type configs struct {
	*AppConfig
	*ProjectConfig
}

func init() {
	envFile := ".env"

	set(nil, &configs{&App, &Project}) // load default value first

	// It should still work even if the .env file does not exist
	if file, err := os.Open(envFile); err == nil {
		set(file, &configs{&App, &Project})
		defer file.Close()
	}
}

func set(file *os.File, structure interface{}) {
	if err := env.Parse(structure); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if err := dotenv.NewDecoder(file).Decode(structure); err != nil {
	}
}
