package feature

import "github.com/rai-project/dlframework"

func AudioType() Option {
	return Type(dlframework.FeatureType_AUDIO)
}

func Audio(e *dlframework.Audio) Option {
	return func(o *dlframework.Feature) {
		AudioType()(o)
		o.Feature = &dlframework.Feature_Audio{
			Audio: e,
		}
	}
}
