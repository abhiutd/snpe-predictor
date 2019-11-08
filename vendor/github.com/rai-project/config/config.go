package config

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// ConfigInterface ...
type ConfigInterface interface {
	ConfigName() string
	SetDefaults()
	Read()
	Wait()
	String() string
	Debug()
}

func setViperConfig(opts *Options) error {
	defer viper.AutomaticEnv() // read in environment variables that match
	defer viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if opts.ConfigString != nil {
		return nil
	}

	if IsValidRemotePrefix(opts.ConfigRemotePath) {
		var provider string
		entry := opts.ConfigRemotePath
		for _, p := range validRemotePrefixes {
			if strings.HasPrefix(entry, p) {
				provider = strings.TrimSuffix(p, "://")
				entry = strings.TrimPrefix(entry, p)
				break
			}
		}
		if !strings.HasPrefix(entry, "http://") &&
			!strings.HasPrefix(entry, "https://") {
			entry = "http://" + entry
		}
		urlParsed, err := url.Parse(entry)
		if err != nil {
			return errors.Errorf("unable to parse remote url %s", entry)
		}
		path := strings.TrimPrefix(urlParsed.Path, "/")
		urlParsed.Scheme = ""
		urlParsed.Path = ""
		urlString := strings.TrimPrefix(urlParsed.String(), "//")
		viper.AddRemoteProvider(provider, urlString, path)
		log.WithField("provider", provider).
			WithField("url", urlString).
			WithField("path", path).
			Info("configuring remote config file provider")
		if strings.HasSuffix(path, ".json") {
			viper.SetConfigType("json")
		} else {
			viper.SetConfigType("yaml")
		}

		return nil
	}

	defer func() {
		for _, pth := range opts.ConfigSearchPaths {
			pth, err := homedir.Expand(pth)
			if err != nil {
				return
			}
			viper.AddConfigPath(pth)
		}
		viper.SetConfigType(opts.ConfigFileType)
	}()

	if com.IsFile(opts.ConfigFileAbsolutePath) {
		log.Debug("Found ", opts.ConfigFileAbsolutePath, " already set. Using ", opts.ConfigFileAbsolutePath, " as the config file.")
		viper.SetConfigFile(opts.ConfigFileAbsolutePath)
		return nil
	}
	if val, ok := os.LookupEnv(opts.ConfigEnvironName); ok {
		pth, _ := homedir.Expand(val)
		log.Debug("Found ", opts.ConfigEnvironName, " in env. Using ", val, " as config file name")
		if com.IsFile(pth) {
			viper.SetConfigFile(pth)
			return nil
		}
		dir, file := path.Split(pth)
		ext := path.Ext(file)
		file = strings.TrimSuffix(file, ext)
		viper.SetConfigName(file)
		viper.AddConfigPath(dir)
		return nil
	}
	if pth, err := filepath.Abs("." + opts.AppName + "_config.yml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ." + opts.AppName + "_config.yml as config file.")
		viper.SetConfigFile(pth)
		return nil
	}
	if pth, err := homedir.Expand("~/." + opts.AppName + "_config.yml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ~/." + opts.AppName + "_config.yml as config file.")
		viper.SetConfigFile(pth)
		return nil
	}
	if pth, err := filepath.Abs("../." + opts.AppName + "_config.yml"); err == nil && com.IsFile(pth) {
		log.Debug("Using ../." + opts.AppName + "_config.yml as config file.")
		viper.SetConfigFile(pth)
		return nil
	}

	log.Info("No fixed configuration file found, searching for a config file with name=", opts.ConfigFileBaseName)
	viper.SetConfigName(opts.ConfigFileBaseName)
	return nil
}

func load(opts *Options) {
	initEnv(opts)
	err := setViperConfig(opts)
	if err != nil {
		log.WithError(err).Error("failed to set viper config")
	}

	if opts.ConfigString != nil {
		configFileName := DefaultAppName + "_config.yml"
		viper.SetConfigFile(configFileName)
		viper.AddConfigPath(".")
		memoryFileSystem := afero.NewMemMapFs()
		file, err := memoryFileSystem.Create(configFileName)
		if err != nil {
			log.WithError(err).Error("Cannot create a memory fs")
		}
		_, err = file.Write([]byte(*opts.ConfigString))
		if err != nil {
			log.WithError(err).Error("cannot write config memory fs")
		}
		defer file.Close()
		viper.SetFs(memoryFileSystem)
	}

	// read configuration
	if IsValidRemotePrefix(opts.ConfigRemotePath) {
		log.Info("reading remote configuration file")
		err = viper.ReadRemoteConfig()
	} else {
		err = viper.ReadInConfig()
	}
	if err != nil {
		log.WithError(err).
			WithField("config_file", viper.ConfigFileUsed()).
			Error("Cannot read in configuration file ")
	}

	// Everything depens on app, so we load it first
	App.SetDefaults()
	App.Read()

	// Load the rest of the configurations
	for _, r := range registry {
		r.SetDefaults()
	}
	for _, r := range registry {
		r.Read()
	}
	// if IsDebug {
	// 	println("read config " + viper.ConfigFileUsed())
	// }
}

// Debug ...
func Debug() {
	log.Debug("Config = ")
	for _, r := range registry {
		r.Debug()
	}
}
