package dlframework

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

func (param *ModelManifest_Type_Parameter) MarshalYAML() (interface{}, error) {
	return cast.ToStringE(param.Value)
}

func (param *ModelManifest_Type_Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	toListString := func(v interface{}) string {
		return fmt.Sprintf("[%s]", strings.Join(cast.ToStringSlice(v), ","))
	}
	var stringSlice []string
	if err := unmarshal(&stringSlice); err == nil {
		param.Value = toListString(stringSlice)
		return nil
	}
	var int32Slice []int32
	if err := unmarshal(&int32Slice); err == nil {
		param.Value = toListString(int32Slice)
		return nil
	}
	var float32Slice []float32
	if err := unmarshal(&float32Slice); err == nil {
		param.Value = toListString(float32Slice)
		return nil
	}

	var int32Value []int32
	if err := unmarshal(&int32Value); err == nil {
		param.Value = cast.ToString(int32Value)
		return nil
	}
	var int64Value []int64
	if err := unmarshal(&int64Value); err == nil {
		param.Value = cast.ToString(int64Value)
		return nil
	}
	var float32Value []float32
	if err := unmarshal(&float32Value); err == nil {
		param.Value = cast.ToString(int32Value)
		return nil
	}
	var float64Value []float64
	if err := unmarshal(&float64Value); err == nil {
		param.Value = cast.ToString(float64Value)
		return nil
	}

	var str string
	if err := unmarshal(&str); err == nil {
		param.Value = str
		return nil
	}
	return errors.New("unable to unmarshal model type parameter")
}
