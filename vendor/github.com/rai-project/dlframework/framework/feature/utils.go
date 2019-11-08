package feature

import (
	"github.com/rai-project/dlframework"
)

func isUnknownType(o *dlframework.Feature) bool {
	return o.Type == dlframework.FeatureType_UNKNOWN
}
