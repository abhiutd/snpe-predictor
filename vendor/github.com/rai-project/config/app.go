package config

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	colorable "github.com/mattn/go-colorable"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/vipertags"
	"github.com/spf13/viper"
)

// DefaultAppDescription ...
const (
	DefaultAppDescription = ""
)

// APP holds common application fields credentials and keys.
type appConfig struct {
	Name        string        `json:"name" config:"app.name"`
	FullName    string        `json:"full_name" config:"app.full_name" default:"rai project"`
	Description string        `json:"description" config:"app.description"`
	License     string        `json:"license" config:"app.license" default:"NCSA"`
	URL         string        `json:"url" config:"app.url" default:"rai-project.com"`
	Secret      string        `json:"-" config:"app.secret"`
	Color       bool          `json:"color" config:"app.color" env:"COLOR"`
	IsDebug     bool          `json:"debug" config:"app.debug" env:"DEBUG"`
	IsVerbose   bool          `json:"verbose" config:"app.verbose" env:"VERBOSE"`
	TempDir     string        `json:"temp_dir" config:"app.tempdir"`
	Version     VersionInfo   `json:"version" config:"-"`
	done        chan struct{} `json:"-" config:"-"`
}

// DefaultAppName ...
var (
	DefaultAppName   = "rai"
	DefaultAppSecret string
	DefaultAppColor  = !color.NoColor
	IsDebug          bool
	IsVerbose        bool

	App = &appConfig{
		Name:        DefaultAppName,
		Description: DefaultAppDescription,
		Version: VersionInfo{
			Version:    "0.2.0",
			BuildDate:  time.Now().String(),
			GitCommit:  "-dirty-",
			GitBranch:  "-dirty-",
			GitState:   "-dirty-",
			GitSummary: "-dirty-",
		},
		done: make(chan struct{}),
	}
)

// ConfigName ...
func (appConfig) ConfigName() string {
	return "App"
}

// SetDefaults ...
func (a *appConfig) SetDefaults() {

	vipertags.SetDefaults(a)
	if a.Name == "" || a.Name == "default" {
		a.Name = DefaultAppName
	}

	viper.SetDefault("app.secret", DefaultAppSecret)
	viper.SetDefault("app.color", DefaultAppColor)
	viper.SetDefault("app.verbose", IsVerbose)
	viper.SetDefault("app.debug", IsDebug)
}

// Read ...
func (a *appConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	if a.Name == "" || a.Name == "default" {
		a.Name = DefaultAppName
	}
	if !viper.IsSet("app.color") {
		a.Color = DefaultAppColor
		viper.Set("app.color", a.Color)
	}
	if a.Secret == "" {
		a.Secret = DefaultAppSecret
	}
	if a.Color == false {
		pp.SetDefaultOutput(colorable.NewNonColorable(pp.GetDefaultOutput()))
	}
	if a.IsDebug || a.IsVerbose {
		pp.WithLineInfo = true
	}
	if a.TempDir == "" {
		a.TempDir, _ = ioutil.TempDir("", a.Name)
	} else {
		expand, err := homedir.Expand(a.TempDir)
		if err != nil {
			return
		}
		a.TempDir = expand
	}
	IsVerbose = a.IsVerbose
	IsDebug = a.IsDebug
}

// Wait ...
func (c appConfig) Wait() {
	<-c.done
}

// String ...
func (a appConfig) String() string {
	return pp.Sprintln(a)
}

// Debug ...
func (a appConfig) Debug() {
	log.Debug("App Config = ", a)
}

func readAppSecret() {
	secretFile, err := homedir.Expand("~/." + App.Name + "_secret")
	if err != nil {
		return
	}
	if !com.IsFile(secretFile) {
		return
	}
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return
	}
	if App.Secret != "" {
		return
	}
	SetAppSecret(strings.TrimSpace(string(b)))
}

// SetAppSecret ...
func SetAppSecret(s string) {
	App.Secret = s
	DefaultAppSecret = s
}

func init() {
	readAppSecret()
}
