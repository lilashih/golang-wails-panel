package helper

import (
	"gbase/src/def"
	"time"
)

func GetNow() string {
	return time.Now().Format(def.YYYY_MM_DD_HH_MM_SS)
}

func GetToday() string {
	return time.Now().Format(def.YYYY_MM_DD)
}
