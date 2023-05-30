package log

import (
	log "github.com/sirupsen/logrus"
	"leblanc.io/open-go-captcha/config"
)

func Initialize(c *config.Config) {
	log.SetLevel(log.ErrorLevel)
}

func Error(message string) {
	log.Error(message)
}