package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/rai-project/config"
)

type Logger struct {
	*logrus.Logger
	mu MutexWrap
}

var (
	debug = false
	std   = New()
)

func UsingHook(s string) bool {
	for _, h := range Config.Hooks {
		if h == s {
			return true
		}
	}
	return false
}

func setupHooks(log *Logger) {
	if UsingHook("stacktrace") && log.Level >= logrus.DebugLevel {
		log.Hooks.Add(StandardStackHook())
	}
	for _, h := range hooks.data {
		log.Hooks.Add(h)
	}
}

func init() {
	config.OnInit(func() {
		config.App.Wait()

		colorMode := config.App.Color || viper.GetBool("app.color")
		formatter := &logrus.TextFormatter{
			DisableColors:    !colorMode,
			ForceColors:      colorMode,
			DisableSorting:   true,
			DisableTimestamp: false,
		}
		logrus.SetFormatter(formatter)
		//logrus.SetOutput(color.Output)
		std.Formatter = formatter

		if config.IsVerbose {
			logrus.SetLevel(logrus.DebugLevel)
			std.Level = logrus.DebugLevel
			Config.Level = "debug"
		} else if config.IsDebug {
			logrus.SetLevel(logrus.DebugLevel)
			std.Level = logrus.DebugLevel
			Config.Level = "debug"
		} else {
			logrus.SetLevel(logrus.InfoLevel)
			std.Level = logrus.InfoLevel
			Config.Level = "info"
		}

		if lvl, err := logrus.ParseLevel(Config.Level); err == nil {
			logrus.SetLevel(lvl)
			std.Level = lvl
		}

		setupHooks(&Logger{Logger: logrus.StandardLogger()})
		setupHooks(std)
	})
}
