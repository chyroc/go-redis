package logger

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry

func init() {
	Log = logrus.NewEntry(logrus.New())
	Log.Level = logrus.DebugLevel
}
