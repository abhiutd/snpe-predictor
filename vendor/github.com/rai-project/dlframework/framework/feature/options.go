package feature

import (
	"github.com/rai-project/dlframework"
)

type Option func(o *dlframework.Feature)

func ID(e string) Option {
	return func(o *dlframework.Feature) {
		o.ID = e
	}
}

func Type(e dlframework.FeatureType) Option {
	return func(o *dlframework.Feature) {
		o.Type = e
	}
}

func Probability(e float32) Option {
	return func(o *dlframework.Feature) {
		o.Probability = e
	}
}

func Metadata(e map[string]string) Option {
	return func(o *dlframework.Feature) {
		o.Metadata = e
	}
}

func AppendMetadata(k, v string) Option {
	return func(o *dlframework.Feature) {
		if o.Metadata == nil {
			o.Metadata = map[string]string{}
		}
		o.Metadata[k] = v
	}
}
