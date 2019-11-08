package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/facebookgo/stack"
	"github.com/fatih/color"
)

// NewHook is the initializer for StackHook{} (implementing logrus.Hook).
// Set levels to callerLevels for which "caller" value may be set, providing a
// single frame of stack. Set levels to stackLevels for which "stack" value may
// be set, providing the full stack (minus logrus).
func NewStackHook(callerLevels []logrus.Level, stackLevels []logrus.Level) StackHook {
	return StackHook{
		CallerLevels: callerLevels,
		StackLevels:  stackLevels,
	}
}

// StackHook is a convenience initializer for StackHook{} with
// default args.
func StandardStackHook() StackHook {
	return StackHook{
		CallerLevels: logrus.AllLevels,
		StackLevels:  []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
	}
}

// StackHook is an implementation of logrus.Hook interface.
type StackHook struct {
	// Set levels to CallerLevels for which "caller" value may be set,
	// providing a single frame of stack.
	CallerLevels []logrus.Level

	// Set levels to StackLevels for which "stack" value may be set,
	// providing the full stack (minus logrus).
	StackLevels []logrus.Level
}

// Levels provides the levels to filter.
func (hook StackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called by logrus when something is logged.
func (hook StackHook) Fire(entry *logrus.Entry) error {
	var skipFrames int
	if len(entry.Data) == 0 {
		// When WithField(s) is not used, we have 8 logrus frames to skip.
		skipFrames = 8
	} else {
		// When WithField(s) is used, we have 6 logrus frames to skip.
		skipFrames = 6
	}

	var frames stack.Stack

	// Get the complete stack track past skipFrames count.
	_frames := stack.Callers(skipFrames)

	// Remove logrus's own frames that seem to appear after the code is through
	// certain hoops. e.g. http handler in a separate package.
	// This is a workaround.
	for _, frame := range _frames {
		if !strings.Contains(frame.File, "github.com/sirupsen/logrus") {
			frames = append(frames, frame)
		}
	}

	if len(frames) > 0 {
		// If we have a frame, we set it to "caller" field for assigned levels.
		for _, level := range hook.CallerLevels {
			if entry.Level == level {
				entry.Data["caller"] = formatFrame(frames[0])
				break
			}
		}

		// Set the available frames to "stack" field.
		for _, level := range hook.StackLevels {
			if entry.Level == level {
				entry.Data["stack"] = frames
				break
			}
		}
	}

	return nil
}

func formatFrame(frame stack.Frame) string {
	green := color.New(color.FgGreen, color.Bold)
	file := green.Sprintf("%v", frame.File)
	line := green.Add(color.Underline).Sprintf("%v", frame.Line)
	name := green.Add(color.Underline).Sprintf("%v", frame.Name)
	return fmt.Sprintf("%s:%s %s", file, line, name)
}
