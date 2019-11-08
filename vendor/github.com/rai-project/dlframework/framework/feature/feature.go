package feature

import (
	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

func New(opts ...Option) *dlframework.Feature {
	feature := &dlframework.Feature{
		ID:       bson.NewObjectId().Hex(),
		Type:     dlframework.FeatureType_UNKNOWN,
		Metadata: map[string]string{},
	}
	for _, o := range opts {
		o(feature)
	}
	return feature
}
