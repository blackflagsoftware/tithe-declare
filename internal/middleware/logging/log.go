package logging

import (
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	"github.com/sirupsen/logrus"
)

var Default = logrus.New()

func init() {
	Default.SetOutput(config.Srv.LogType)
	Default.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}
