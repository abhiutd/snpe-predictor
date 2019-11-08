package feature

import "github.com/rai-project/dlframework"

func RegionType() Option {
	return Type(dlframework.FeatureType_REGION)
}

func Region(e *dlframework.Region) Option {
	return func(o *dlframework.Feature) {
		RegionType()(o)
		o.Feature = &dlframework.Feature_Region{
			Region: e,
		}
	}
}
