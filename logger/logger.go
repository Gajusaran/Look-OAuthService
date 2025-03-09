package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}
