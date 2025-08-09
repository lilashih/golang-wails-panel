package helper

import (
	"gbase/src/core/config"
)

func IsRelease() bool {
	return config.App.Mode == "release"
}

func IsDebug() bool {
	return config.App.Mode == "debug"
}

func IsTest() bool {
	return config.App.Mode == "test"
}
