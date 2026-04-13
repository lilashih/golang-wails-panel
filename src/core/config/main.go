package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/golobby/dotenv"
)

var App AppConfig
var Project ProjectConfig
var Logger LoggerConfig

type configs struct {
	*AppConfig
	*ProjectConfig
	*LoggerConfig
}

func init() {
	envFile := ".env"

	set(nil, &configs{&App, &Project, &Logger}) // load default value first

	// It should still work even if the .env file does not exist
	if file, err := openEnvFile(envFile); err == nil {
		set(file, &configs{&App, &Project, &Logger})
		defer file.Close()
	}
}

func openEnvFile(envFile string) (*os.File, error) {
	file, err := os.Open(envFile)
	if err == nil {
		return file, nil
	}

	if App.Mode != "release" {
		return nil, err
	}

	if appImagePath := os.Getenv("APPIMAGE"); appImagePath != "" {
		return os.Open(filepath.Join(filepath.Dir(appImagePath), envFile))
	}

	executablePath, executableErr := os.Executable()
	if executableErr != nil {
		return nil, err
	}

	return os.Open(filepath.Join(filepath.Dir(executablePath), envFile))
}

func set(file *os.File, structure interface{}) {
	if err := env.Parse(structure); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if err := dotenv.NewDecoder(file).Decode(structure); err != nil {
	}
}
