package logger

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/spf13/viper"
)

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

// Creates a new logrusger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logrusger instance. You can also just
// instantiate your own:
//
//    var logrus = &Logger{
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrusrus.DebugLevel,
//    }
//
// It's recommended to make this a global instance called `logrus`.
func New() *Logger {
	l := logrus.New()
	colorMode := config.App.Color || viper.GetBool("app.color")
	formatter := &logrus.TextFormatter{
		DisableColors:    !colorMode,
		ForceColors:      colorMode,
		DisableSorting:   true,
		DisableTimestamp: false,
	}
	l.Formatter = formatter
	//l.Out = color.Output
	//if !viper.GetBool("app.color") {
	//	l.Out = colorable.NewNonColorable(os.Stdout)
	//}
	if config.IsVerbose {
		l.Level = logrus.DebugLevel
	} else {
		l.Level = logrus.WarnLevel
	}
	logrusger := &Logger{
		Logger: l,
	}
	setupHooks(logrusger)
	return logrusger
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

func StandardLogger() *Logger {
	return std
}

// SetOutput sets the standard logrusger output.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Out = out
}

// SetFormatter sets the standard logrusger formatter.
func SetFormatter(formatter logrus.Formatter) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Formatter = formatter
}

// SetLevel sets the standard logrusger level.
func SetLevel(level logrus.Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Level = level
}

// GetLevel returns the standard logrusger level.
func GetLevel() logrus.Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.Level
}

// AddHook adds a hook to the standard logrusger hooks.
func AddHook(hook logrus.Hook) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Hooks.Add(hook)
}

// WithError creates an entry from the standard logrusger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return std.WithField(logrus.ErrorKey, err)
}

// WithField creates an entry from the standard logrusger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't logrus until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logrusger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't logrus until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return std.WithFields(fields)
}

// Debug logruss a message at level Debug on the standard logrusger.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Print logruss a message at level Info on the standard logrusger.
func Print(args ...interface{}) {
	std.Print(args...)
}

// Info logruss a message at level Info on the standard logrusger.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logruss a message at level Warn on the standard logrusger.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warning logruss a message at level Warn on the standard logrusger.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Error logruss a message at level Error on the standard logrusger.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Panic logruss a message at level Panic on the standard logrusger.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Fatal logruss a message at level Fatal on the standard logrusger.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Debugf logruss a message at level Debug on the standard logrusger.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Printf logruss a message at level Info on the standard logrusger.
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Infof logruss a message at level Info on the standard logrusger.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logruss a message at level Warn on the standard logrusger.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Warningf logruss a message at level Warn on the standard logrusger.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Errorf logruss a message at level Error on the standard logrusger.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Panicf logruss a message at level Panic on the standard logrusger.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Fatalf logruss a message at level Fatal on the standard logrusger.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Debugln logruss a message at level Debug on the standard logrusger.
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Println logruss a message at level Info on the standard logrusger.
func Println(args ...interface{}) {
	std.Println(args...)
}

// Infoln logruss a message at level Info on the standard logrusger.
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Warnln logruss a message at level Warn on the standard logrusger.
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

// Warningln logruss a message at level Warn on the standard logrusger.
func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

// Errorln logruss a message at level Error on the standard logrusger.
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Panicln logruss a message at level Panic on the standard logrusger.
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

// Fatalln logruss a message at level Fatal on the standard logrusger.
func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}
