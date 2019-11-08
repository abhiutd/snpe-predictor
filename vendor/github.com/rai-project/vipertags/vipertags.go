package vipertags

import (
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"os"

	"github.com/sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func setField(field *structs.Field, val interface{}) {
	switch field.Value().(type) {
	case bool:
		field.Set(cast.ToBool(val))
	case string:
		field.Set(cast.ToString(val))
	case int64:
		field.Set(cast.ToInt64(val))
	case int32, int16, int8, int:
		field.Set(cast.ToInt(val))
	case float64, float32:
		field.Set(cast.ToFloat64(val))
	case time.Time:
		field.Set(cast.ToTime(val))
	case time.Duration:
		field.Set(cast.ToDuration(val))
	case map[string]string:
		field.Set(cast.ToStringMapString(val))
	case map[string][]string:
		field.Set(cast.ToStringMapStringSlice(val))
	case map[string]bool:
		field.Set(cast.ToStringMapBool(val))
	case map[string]interface{}:
		field.Set(cast.ToStringMap(val))
	case []int:
		field.Set(cast.ToIntSlice(val))
	case []string:
		field.Set(cast.ToStringSlice(val))
	case []interface{}:
		field.Set(cast.ToSlice(val))
	default:
		field.Set(val)
	}
}

func buildDefault(st0 interface{}, prefix0 string) interface{} {
	st := structs.New(st0)
	for _, field := range st.Fields() {

		configTagValue := field.Tag("config")

		if configTagValue == "-" {
			continue
		}

		defaultTagValue := field.Tag("default")

		prefix := prefix0
		if configTagValue != "" {
			prefix = prefix + configTagValue
		}

		if field.Kind() == reflect.Struct {
			buildDefault(field.Value(), prefix)
			continue
		}

		if defaultTagValue != "" {
			setField(field, defaultTagValue)
		}
	}
	return st0
}

func buildConfiguration(st0 interface{}, prefix0 string) interface{} {
	st := structs.New(st0)
	for _, field := range st.Fields() {

		configTagValue := field.Tag("config")

		if configTagValue == "-" {
			continue
		}

		envTagValue := field.Tag("env")

		prefix := prefix0
		if configTagValue != "" {
			prefix = prefix + configTagValue
		}

		if field.Kind() == reflect.Struct {
			buildConfiguration(field.Value(), prefix)
			continue
		}

		if envTagValue != "" && configTagValue != "" {
			viper.BindEnv(configTagValue, envTagValue)
		}
		if envTagValue != "" && configTagValue == "" {
			if e := os.Getenv(envTagValue); e != "" {
				setField(field, e)
			}
		}
		if configTagValue != "" && viper.IsSet(configTagValue) {
			setField(field, viper.Get(configTagValue))
		}
	}
	return st0
}

func SetDefaults(class interface{}) {
	buildDefault(class, "")
}

func Fill(class interface{}) {
	buildConfiguration(class, "")
}

func Setup(fileType string, prefix string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("conf")
	viper.AddConfigPath("config")
	viper.SetConfigType(fileType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func FromFile(filename string, prefix string) {
	Setup(strings.Replace(filepath.Ext(filename), ".", "", 1), prefix)
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).
			Fatal("Cannot find configuration file.")
	}
}
