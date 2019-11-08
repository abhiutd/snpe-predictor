package config

import (
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// Options ...
type Options struct {
	AppName                string
	AppSecret              string
	ConfigSearchPaths      []string
	ConfigEnvironName      string
	ConfigFileBaseName     string
	ConfigFileType         string
	ConfigFileAbsolutePath string
	ConfigString           *string
	ConfigRemotePath       string
	IsVerbose              bool
	IsDebug                bool
}

// Option ...
type Option func(*Options)

// NewOptions ...
func NewOptions() *Options {
	isVerbose, isDebug := modeInfo()
	return &Options{
		AppName:            DefaultAppName,
		AppSecret:          DefaultAppSecret,
		ConfigSearchPaths:  []string{"$HOME", "..", "../..", "."},
		ConfigEnvironName:  strings.ToUpper(DefaultAppName) + "_CONFIG_FILE",
		ConfigFileBaseName: "." + strings.ToLower(DefaultAppName) + "_config",
		ConfigFileType:     "yaml",
		ConfigString:       nil,
		IsDebug:            isDebug,
		IsVerbose:          isVerbose,
	}
}

// AppName ...
func AppName(s string) Option {
	return func(opts *Options) {
		DefaultAppName = s
		opts.AppName = s
		opts.ConfigFileBaseName = "." + strings.ToLower(DefaultAppName) + "_config"
	}
}

// AppSecret ...
func AppSecret(s string) Option {
	return func(opts *Options) {
		DefaultAppSecret = s
		opts.AppSecret = s
	}
}

// ConfigSearchPaths ...
func ConfigSearchPaths(s []string) Option {
	return func(opts *Options) {
		opts.ConfigSearchPaths = s
	}
}

// ConfigEnvironName ...
func ConfigEnvironName(s string) Option {
	return func(opts *Options) {
		opts.ConfigEnvironName = s
	}
}

// ConfigFileBaseName ...
func ConfigFileBaseName(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileBaseName = s
	}
}

// ConfigFileType ...
func ConfigFileType(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileType = s
	}
}

// ConfigRemotePath ...
func ConfigRemotePath(s string) Option {
	return func(opts *Options) {
		opts.ConfigRemotePath = s
	}
}

// ConfigFileAbsolutePath ...
func ConfigFileAbsolutePath(s string) Option {
	return func(opts *Options) {
		opts.ConfigFileAbsolutePath = s
	}
}

// ConfigString ...
func ConfigString(s string) Option {
	return func(opts *Options) {
		opts.ConfigString = &s
	}
}

// VerboseMode ...
func VerboseMode(b bool) Option {
	return func(opts *Options) {
		opts.IsVerbose = b
		IsVerbose = b
		App.IsVerbose = b
		viper.Set("app.verbose", b)
	}
}

// DebugMode ...
func DebugMode(b bool) Option {
	return func(opts *Options) {
		opts.IsDebug = b
		IsDebug = b
		App.IsDebug = b
		viper.Set("app.debug", b)
	}
}

// ColorMode ...
func ColorMode(b bool) Option {
	return func(opts *Options) {
		DefaultAppColor = b
		App.Color = b
		color.NoColor = !b
		viper.Set("app.color", b)
	}
}
